// Package logic 兑换码业务(移植自 tianbi redeemcodeer)。
// 与 tianbi 差异: my_service 用户当前只有金币(users.balance), card_type 仅启用 1=金币;
// 防超发用「used_times < total_times 条件更新 + 影响行数判定」原子占额;
// 防同一用户重复兑换用查询 + 唯一约束(site_id,user_id,code)双保险。
package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/redeemcode/service"
	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const (
	rcSiteId      = 1                                  // 单站点样板
	rcCardGold    = 1                                  // card_type: 金币
	rcCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 自动生成码字符集(去易混淆字符)
	rcCodeLen     = 12
)

type sRedeemCode struct{}

func New() service.IRedeemCode { return &sRedeemCode{} }

func goldDesc(value int) string { return fmt.Sprintf("兑换金币x%d", value) }

// Use 使用兑换码(移植自 tianbi UseRedeemCode)。
func (s *sRedeemCode) Use(ctx context.Context, userId int64, code string) (*service.UseResult, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, gerror.New("兑换码不能为空")
	}
	// 1. 该用户是否已用过此码
	used, err := g.Model("redeem_code_record").Ctx(ctx).
		Where("site_id", rcSiteId).Where("user_id", userId).Where("code", code).Count()
	if err != nil {
		return nil, err
	}
	if used > 0 {
		return nil, gerror.New("您已经使用过该兑换码了")
	}
	// 2. 查码并校验
	var rc *entity.RedeemCode
	if err := g.Model("redeem_code").Ctx(ctx).
		Where("site_id", rcSiteId).Where("code", code).Scan(&rc); err != nil {
		return nil, err
	}
	if rc == nil {
		return nil, gerror.New("兑换码不存在")
	}
	if rc.Status != 1 {
		return nil, gerror.New("兑换码已失效")
	}
	if rc.ExpiredAt != nil && rc.ExpiredAt.Before(gtime.Now()) {
		return nil, gerror.New("兑换码已过期")
	}
	if rc.UsedTimes >= rc.TotalTimes {
		return nil, gerror.New("兑换码兑换次数已用完")
	}
	if rc.CardType != rcCardGold {
		return nil, gerror.New("该兑换码类型暂不支持")
	}
	// 3. 事务: 原子占额 → 发金币 → 记录
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("redeem_code").Ctx(ctx).
			Where("id", rc.Id).Where("used_times < total_times").
			Data(g.Map{
				"used_times": &gdb.Counter{Field: "used_times", Value: 1},
				"updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 { // 并发下被别人抢完
			return gerror.New("兑换码兑换次数已用完")
		}
		// 次数用完自动禁用
		if _, err := tx.Model("redeem_code").Ctx(ctx).
			Where("id", rc.Id).Where("used_times >= total_times").
			Data(g.Map{"status": 0}).Update(); err != nil {
			return err
		}
		if rc.Value > 0 {
			if err := balance.Add(ctx, tx, userId, float64(rc.Value), balance.SceneRedeemCode, fmt.Sprintf("redeem_code:%d", rc.Id), goldDesc(rc.Value)); err != nil {
				return err
			}
		}
		// 记录(唯一约束兜底并发下同用户重复提交)
		if _, err := tx.Model("redeem_code_record").Ctx(ctx).Data(g.Map{
			"site_id": rcSiteId, "user_id": userId, "code_id": rc.Id,
			"code": rc.Code, "name": rc.Name, "card_type": rc.CardType, "value": rc.Value,
		}).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &service.UseResult{Desc: goldDesc(rc.Value)}, nil
}

// MyRecords 我的兑换记录(移植自 tianbi GetRedeemCodeRecords)。
func (s *sRedeemCode) MyRecords(ctx context.Context, userId int64, page, size int) ([]service.MyRecordDTO, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var list []*entity.RedeemCodeRecord
	if err := g.Model("redeem_code_record").Ctx(ctx).
		Where("site_id", rcSiteId).Where("user_id", userId).
		OrderDesc("created_at").Page(page, size).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]service.MyRecordDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, service.MyRecordDTO{
			Code: r.Code, Desc: goldDesc(r.Value), ActivedAt: created,
		})
	}
	return out, nil
}

func (s *sRedeemCode) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("redeem_code").Ctx(ctx).Where("site_id", rcSiteId)
	if f.Status >= 0 { // -1=全部
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(code ILIKE ? OR name ILIKE ?)", kw, kw)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RedeemCode
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, &service.ItemDTO{
			Id: r.Id, Name: r.Name, Code: r.Code, CardType: r.CardType, Value: r.Value,
			TotalTimes: r.TotalTimes, UsedTimes: r.UsedTimes, Status: r.Status,
			ExpiredAt: timeStr(r.ExpiredAt), CreatedAt: timeStr(r.CreatedAt),
		})
	}
	return out, total, nil
}

func (s *sRedeemCode) Create(ctx context.Context, in service.CreateInput) (int64, string, error) {
	if in.Value <= 0 {
		return 0, "", gerror.New("金币数需大于0")
	}
	if in.TotalTimes <= 0 {
		return 0, "", gerror.New("次数至少为1")
	}
	exp := gtime.NewFromStr(in.ExpiredAt)
	if exp == nil || exp.IsZero() {
		return 0, "", gerror.New("过期时间格式非法, 如: 2027-12-31 23:59:59")
	}
	if exp.Before(gtime.Now()) {
		return 0, "", gerror.New("过期时间不能早于现在")
	}
	code := strings.ToUpper(strings.TrimSpace(in.Code))
	if code == "" {
		code = grand.Str(rcCodeCharset, rcCodeLen)
	}
	cnt, err := g.Model("redeem_code").Ctx(ctx).
		Where("site_id", rcSiteId).Where("code", code).Count()
	if err != nil {
		return 0, "", err
	}
	if cnt > 0 {
		return 0, "", gerror.New("该兑换码已存在")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("redeem_code").Ctx(ctx).Data(g.Map{
		"site_id": rcSiteId, "name": in.Name, "code": code, "card_type": rcCardGold,
		"value": in.Value, "total_times": in.TotalTimes, "status": in.Status,
		"expired_at": exp,
	}).InsertAndGetId()
	if err != nil {
		return 0, "", err
	}
	return id, code, nil
}

func (s *sRedeemCode) Update(ctx context.Context, in service.UpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"updated_at": gtime.Now()}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Value > 0 {
		data["value"] = in.Value
	}
	if in.TotalTimes > 0 {
		data["total_times"] = in.TotalTimes
	}
	if in.ExpiredAt != "" {
		exp := gtime.NewFromStr(in.ExpiredAt)
		if exp == nil || exp.IsZero() {
			return gerror.New("过期时间格式非法")
		}
		data["expired_at"] = exp
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("redeem_code").Ctx(ctx).
		Where("site_id", rcSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sRedeemCode) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("redeem_code").Ctx(ctx).
		Where("site_id", rcSiteId).Where("id", id).Delete()
	return err
}

func (s *sRedeemCode) Records(ctx context.Context, f service.RecordFilter) ([]*service.RecordDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("redeem_code_record").Ctx(ctx).Where("site_id", rcSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Code != "" {
		m = m.Where("code", strings.ToUpper(strings.TrimSpace(f.Code)))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RedeemCodeRecord
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.RecordDTO, 0, len(list))
	for _, r := range list {
		out = append(out, &service.RecordDTO{
			Id: r.Id, UserId: r.UserId, CodeId: r.CodeId, Code: r.Code, Name: r.Name,
			CardType: r.CardType, Value: r.Value, CreatedAt: timeStr(r.CreatedAt),
		})
	}
	return out, total, nil
}

func timeStr(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}
