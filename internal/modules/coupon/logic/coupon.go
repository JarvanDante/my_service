// Package logic 优惠券业务(移植自 tianbi coupons/coupons_user)。
//
// 修掉的 tianbi 老问题:
//  1. 自动选券时的过期过滤写成 `{"gte": now}`(少了 `$`), 过期券会被选中 →
//     这里把「状态=未使用 + 未过期 + 场景匹配 + 门槛达标」四个条件全写进 SQL;
//  2. 发券没有总量控制 → 这里用 `WHERE total = -1 OR issued < total` 条件递增, 防超发;
//  3. 核销没有并发保护 → 这里 `WHERE id=? AND user_id=? AND status=1` 条件更新 + RowsAffected 判定,
//     并要求调用方传入事务, 与扣款/下单同生共死。
package logic

import (
	"context"
	"math"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/coupon/service"
)

const cpSiteId = 1

type sCoupon struct{}

func New() service.ICoupon { return &sCoupon{} }

func round2(f float64) float64 { return math.Round(f*100) / 100 }

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

// deductOf 计算某张券对给定订单金额的抵扣额(服务端唯一算法, 客户端不参与)。
// 抵用券: 直接减面额; 折扣券: amount*(100-discount)/100, 受 max_deduct 封顶。
// 结果被 clamp 在 [0, amount], 保证不会出现负数应付。
func deductOf(c *entity.UserCoupon, amount float64) float64 {
	if amount < c.Threshold {
		return 0
	}
	var d float64
	switch c.Type {
	case service.TypeCash:
		d = c.FaceValue
	case service.TypeDiscount:
		disc := c.Discount
		if disc <= 0 || disc >= 100 {
			return 0
		}
		d = amount * float64(100-disc) / 100
		if c.MaxDeduct > 0 && d > c.MaxDeduct {
			d = c.MaxDeduct
		}
	default:
		return 0
	}
	if d > amount {
		d = amount
	}
	if d < 0 {
		d = 0
	}
	return round2(d)
}

func toTplDTO(r *entity.CouponTpl) *service.TplDTO {
	return &service.TplDTO{
		Id: r.Id, Name: r.Name, Type: r.Type, Scene: r.Scene, FaceValue: r.FaceValue,
		Discount: r.Discount, Threshold: r.Threshold, MaxDeduct: r.MaxDeduct,
		Total: r.Total, Issued: r.Issued, PerLimit: r.PerLimit, ExpireDay: r.ExpireDay,
		Status: r.Status, CreatedAt: fmtTime(r.CreatedAt),
	}
}

func toUserDTO(r *entity.UserCoupon) *service.UserCouponDTO {
	return &service.UserCouponDTO{
		Id: r.Id, UserId: r.UserId, TplId: r.TplId, Name: r.Name, Type: r.Type,
		Scene: r.Scene, FaceValue: r.FaceValue, Discount: r.Discount,
		Threshold: r.Threshold, MaxDeduct: r.MaxDeduct, Status: r.Status, RefId: r.RefId,
		ExpireAt: fmtTime(r.ExpireAt), UsedAt: fmtTime(r.UsedAt), CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sCoupon) Tpls(ctx context.Context, userId int64) ([]*service.TplDTO, error) {
	var list []*entity.CouponTpl
	if err := g.Model("coupon_tpl").Ctx(ctx).
		Where("site_id", cpSiteId).Where("status", 1).
		Where("total = -1 OR issued < total").
		OrderDesc("id").Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.TplDTO, 0, len(list))
	for _, r := range list {
		d := toTplDTO(r)
		if userId > 0 && r.PerLimit > 0 {
			n, err := g.Model("user_coupon").Ctx(ctx).
				Where("site_id", cpSiteId).Where("user_id", userId).Where("tpl_id", r.Id).Count()
			if err != nil {
				return nil, err
			}
			d.Received = n >= r.PerLimit
		}
		out = append(out, d)
	}
	return out, nil
}

// issue 发一张券(领取 / 后台发放共用)。tx 可为 nil, 为 nil 时用普通连接。
func (s *sCoupon) issue(ctx context.Context, userId, tplId int64) (int64, error) {
	var tpl *entity.CouponTpl
	if err := g.Model("coupon_tpl").Ctx(ctx).
		Where("site_id", cpSiteId).Where("id", tplId).Where("status", 1).Scan(&tpl); err != nil {
		return 0, err
	}
	if tpl == nil {
		return 0, gerror.New("券不存在或已停用")
	}
	if tpl.PerLimit > 0 {
		n, err := g.Model("user_coupon").Ctx(ctx).
			Where("site_id", cpSiteId).Where("user_id", userId).Where("tpl_id", tplId).Count()
		if err != nil {
			return 0, err
		}
		if n >= tpl.PerLimit {
			return 0, gerror.New("已达该券领取上限")
		}
	}
	expireDay := tpl.ExpireDay
	if expireDay <= 0 {
		expireDay = 7
	}
	var newId int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 总量条件递增: -1 不限量; 否则必须 issued < total
		res, err := tx.Model("coupon_tpl").Ctx(ctx).
			Where("id", tplId).Where("total = -1 OR issued < total").
			Data(g.Map{
				"issued":     &gdb.Counter{Field: "issued", Value: 1},
				"updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("该券已领完")
		}
		id, err := tx.Model("user_coupon").Ctx(ctx).Data(g.Map{
			"site_id": cpSiteId, "user_id": userId, "tpl_id": tplId, "name": tpl.Name,
			"type": tpl.Type, "scene": tpl.Scene, "face_value": tpl.FaceValue,
			"discount": tpl.Discount, "threshold": tpl.Threshold, "max_deduct": tpl.MaxDeduct,
			"status": service.StatusUnused, "expire_at": gtime.Now().AddDate(0, 0, expireDay),
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newId = id
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *sCoupon) Receive(ctx context.Context, userId, tplId int64) (int64, error) {
	return s.issue(ctx, userId, tplId)
}

// expireStale 把已到期的未使用券刷成"已过期"(懒过期, 不用定时任务)。
func (s *sCoupon) expireStale(ctx context.Context, userId int64) error {
	m := g.Model("user_coupon").Ctx(ctx).
		Where("site_id", cpSiteId).Where("status", service.StatusUnused).
		Where("expire_at < ?", gtime.Now())
	if userId > 0 {
		m = m.Where("user_id", userId)
	}
	_, err := m.Data(g.Map{"status": service.StatusExpire}).Update()
	return err
}

func (s *sCoupon) My(ctx context.Context, userId int64, status, page, size int) ([]*service.UserCouponDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if err := s.expireStale(ctx, userId); err != nil {
		return nil, 0, err
	}
	m := g.Model("user_coupon").Ctx(ctx).Where("site_id", cpSiteId).Where("user_id", userId)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCoupon
	if err := m.Clone().OrderAsc("status").OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.UserCouponDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toUserDTO(r))
	}
	return out, total, nil
}

func (s *sCoupon) Available(ctx context.Context, userId int64, scene int, amount float64) ([]*service.UserCouponDTO, int64, float64, error) {
	if amount <= 0 {
		return nil, 0, 0, gerror.New("订单金额需大于0")
	}
	if err := s.expireStale(ctx, userId); err != nil {
		return nil, 0, 0, err
	}
	m := g.Model("user_coupon").Ctx(ctx).
		Where("site_id", cpSiteId).Where("user_id", userId).
		Where("status", service.StatusUnused).
		Where("expire_at > ?", gtime.Now()).
		Where("threshold <= ?", amount)
	if scene > 0 && scene != service.SceneAll {
		m = m.WhereIn("scene", g.Slice{scene, service.SceneAll})
	}
	var list []*entity.UserCoupon
	if err := m.OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, 0, err
	}
	out := make([]*service.UserCouponDTO, 0, len(list))
	var bestId int64
	var bestDeduct float64
	for _, r := range list {
		d := deductOf(r, amount)
		if d <= 0 {
			continue
		}
		dto := toUserDTO(r)
		dto.Deduct = d
		out = append(out, dto)
		if d > bestDeduct {
			bestDeduct, bestId = d, r.Id
		}
	}
	return out, bestId, bestDeduct, nil
}

func (s *sCoupon) UseInTx(ctx context.Context, tx gdb.TX, userId, couponId int64,
	scene int, amount float64, refId string) (float64, error) {
	if couponId <= 0 {
		return 0, nil // 不用券
	}
	var c *entity.UserCoupon
	if err := tx.Model("user_coupon").Ctx(ctx).
		Where("site_id", cpSiteId).Where("id", couponId).Where("user_id", userId).
		Scan(&c); err != nil {
		return 0, err
	}
	if c == nil {
		return 0, gerror.New("优惠券不存在")
	}
	if c.Status != service.StatusUnused {
		return 0, gerror.New("优惠券不可用")
	}
	if c.ExpireAt != nil && c.ExpireAt.Before(gtime.Now()) {
		return 0, gerror.New("优惠券已过期")
	}
	if scene > 0 && c.Scene != service.SceneAll && c.Scene != scene {
		return 0, gerror.New("优惠券不适用于该场景")
	}
	d := deductOf(c, amount)
	if d <= 0 {
		return 0, gerror.New("订单金额未达该券使用门槛")
	}
	res, err := tx.Model("user_coupon").Ctx(ctx).
		Where("id", couponId).Where("user_id", userId).Where("status", service.StatusUnused).
		Data(g.Map{
			"status": service.StatusUsed, "used_at": gtime.Now(), "ref_id": refId,
		}).Update()
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, gerror.New("优惠券已被使用")
	}
	return d, nil
}

// ---------------- 后台 ----------------

func (s *sCoupon) List(ctx context.Context, f service.ListFilter) ([]*service.TplDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("coupon_tpl").Ctx(ctx).Where("site_id", cpSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.CouponTpl
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.TplDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toTplDTO(r))
	}
	return out, total, nil
}

func normalizeTpl(in *service.TplInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return gerror.New("券名不能为空")
	}
	if in.Type != service.TypeCash && in.Type != service.TypeDiscount {
		in.Type = service.TypeCash
	}
	if in.Scene <= 0 || in.Scene > service.SceneAll {
		in.Scene = service.SceneAll
	}
	if in.Type == service.TypeCash && in.FaceValue <= 0 {
		return gerror.New("抵用券面额需大于0")
	}
	if in.Type == service.TypeDiscount && (in.Discount <= 0 || in.Discount >= 100) {
		return gerror.New("折扣需在 1~99 之间(85 表示85折)")
	}
	if in.Total == 0 {
		in.Total = -1
	}
	if in.ExpireDay <= 0 {
		in.ExpireDay = 7
	}
	if in.Discount <= 0 {
		in.Discount = 100
	}
	return nil
}

func (s *sCoupon) Create(ctx context.Context, in service.TplInput) (int64, error) {
	if err := normalizeTpl(&in); err != nil {
		return 0, err
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	return g.Model("coupon_tpl").Ctx(ctx).Data(g.Map{
		"site_id": cpSiteId, "name": in.Name, "type": in.Type, "scene": in.Scene,
		"face_value": in.FaceValue, "discount": in.Discount, "threshold": in.Threshold,
		"max_deduct": in.MaxDeduct, "total": in.Total, "per_limit": in.PerLimit,
		"expire_day": in.ExpireDay, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sCoupon) Update(ctx context.Context, in service.TplInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{
		"threshold": in.Threshold, "max_deduct": in.MaxDeduct,
		"per_limit": in.PerLimit, "updated_at": gtime.Now(),
	}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Type == service.TypeCash || in.Type == service.TypeDiscount {
		data["type"] = in.Type
	}
	if in.Scene > 0 && in.Scene <= service.SceneAll {
		data["scene"] = in.Scene
	}
	if in.FaceValue > 0 {
		data["face_value"] = in.FaceValue
	}
	if in.Discount > 0 && in.Discount <= 100 {
		data["discount"] = in.Discount
	}
	if in.Total != 0 {
		data["total"] = in.Total
	}
	if in.ExpireDay > 0 {
		data["expire_day"] = in.ExpireDay
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("coupon_tpl").Ctx(ctx).
		Where("site_id", cpSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sCoupon) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	n, err := g.Model("user_coupon").Ctx(ctx).
		Where("site_id", cpSiteId).Where("tpl_id", id).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return gerror.New("该券已发放给用户, 不能删除, 请改为停用")
	}
	_, err = g.Model("coupon_tpl").Ctx(ctx).
		Where("site_id", cpSiteId).Where("id", id).Delete()
	return err
}

func (s *sCoupon) Grant(ctx context.Context, tplId int64, userIds []int64) (int, int, []string) {
	var ok, fail int
	errs := make([]string, 0)
	for _, uid := range userIds {
		if uid <= 0 {
			fail++
			errs = append(errs, "用户ID非法")
			continue
		}
		if _, err := s.issue(ctx, uid, tplId); err != nil {
			fail++
			errs = append(errs, "user "+gconv.String(uid)+": "+err.Error())
			continue
		}
		ok++
	}
	return ok, fail, errs
}

func (s *sCoupon) Users(ctx context.Context, f service.ListFilter) ([]*service.UserCouponDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	if err := s.expireStale(ctx, 0); err != nil {
		return nil, 0, err
	}
	m := g.Model("user_coupon").Ctx(ctx).Where("site_id", cpSiteId)
	if f.TplId > 0 {
		m = m.Where("tpl_id", f.TplId)
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Status > 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCoupon
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.UserCouponDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toUserDTO(r))
	}
	return out, total, nil
}
