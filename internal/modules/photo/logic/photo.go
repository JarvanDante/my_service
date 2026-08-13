// Package logic 图集业务。
//
// 关键取舍:
//   - 图集没有章节, 一条记录就是一整套图, 付费是「整套购买」, 解锁判定复用 shared/paywall;
//   - 漫画的"前 N 章免费"在这里退化为"前 N 张免费": 未解锁时详情只返回 pics[:free_count],
//     截断必须在服务端做 —— 若把全部 url 下发再让客户端遮罩, 付费墙形同虚设;
//   - 图片存 jsonb 整体读写, 另冗余 pic_count 供列表页展示, 避免列表里解 jsonb;
//   - 购买价取库里的 price, 绝不接受客户端传值。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/photo/service"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const phSiteId = 1

type sPhoto struct{}

func New() service.IPhoto { return &sPhoto{} }

func decodeTags(raw string) []string {
	out := []string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func encodeJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodePics(raw string) []service.PicDTO {
	out := []service.PicDTO{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func toDTO(r *entity.PhotoAlbum) *service.PhotoDTO {
	return &service.PhotoDTO{
		Id: r.Id, Title: r.Title, Cover: r.Cover, Intro: r.Intro, Category: r.Category,
		Tags: decodeTags(r.Tags), IsVip: r.IsVip, Price: r.Price, FreeCount: r.FreeCount,
		Pics: decodePics(r.Pics), PicCount: r.PicCount, ViewCount: r.ViewCount,
		BuyCount: r.BuyCount, LikeCount: r.LikeCount, Rank: r.Rank, Status: r.Status,
		PublishId: r.PublishId, CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sPhoto) query(ctx context.Context, f service.ListFilter) ([]*service.PhotoDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	base := g.Model("photo_album").Ctx(ctx).Where("site_id", phSiteId)
	if f.Status >= 0 {
		base = base.Where("status", f.Status)
	}
	if f.Category != "" {
		base = base.Where("category", f.Category)
	}
	if f.Tag != "" {
		base = base.Where("tags @> ?::jsonb", encodeJSON([]string{f.Tag}))
	}
	if f.Keyword != "" {
		base = base.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	m := base.Clone()
	switch f.Sort {
	case 1:
		m = m.OrderDesc("view_count")
	case 2:
		m = m.OrderDesc("id")
	case 3:
		m = m.OrderDesc("like_count")
	default:
		m = m.OrderDesc("rank")
	}
	var list []*entity.PhotoAlbum
	if err := m.OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.PhotoDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sPhoto) FrontList(ctx context.Context, userId int64, f service.ListFilter) ([]*service.PhotoDTO, int, error) {
	f.Status = entity.ContentStatusOnline
	list, total, err := s.query(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	for _, d := range list {
		d.Pics = nil // 列表页不下发图片: 既省流量, 也避免未购用户从列表接口套出全图
	}
	if userId > 0 && len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, d := range list {
			ids = append(ids, d.Id)
		}
		bought, err := paywall.PurchasedSet(ctx, userId, paywall.MediaPhoto, ids)
		if err != nil {
			return nil, 0, err
		}
		for _, d := range list {
			d.IsBuy = bought[d.Id]
		}
	}
	return list, total, nil
}

// find 取图集; onlineOnly=true 时只取上架的(前台用)。
func (s *sPhoto) find(ctx context.Context, id int64, onlineOnly bool) (*entity.PhotoAlbum, error) {
	m := g.Model("photo_album").Ctx(ctx).Where("site_id", phSiteId).Where("id", id)
	if onlineOnly {
		m = m.Where("status", entity.ContentStatusOnline)
	}
	var r *entity.PhotoAlbum
	if err := m.Scan(&r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, gerror.New("图集不存在或已下架")
	}
	return r, nil
}

func (s *sPhoto) Detail(ctx context.Context, userId, id int64) (*service.DetailDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return nil, err
	}
	acc, err := paywall.Check(ctx, userId, paywall.MediaPhoto, r.Id, r.IsVip == 1, r.Price)
	if err != nil {
		return nil, err
	}
	// 观看数 +1(失败不影响主流程)
	_, _ = g.Model("photo_album").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"view_count": &gdb.Counter{Field: "view_count", Value: 1}}).Update()
	d := toDTO(r)
	d.IsBuy = acc.IsBuy
	d.ViewCount++
	total := len(d.Pics)
	if !acc.Playable {
		// 未解锁: 只放行前 free_count 张。free_count 可能被配成 0 或超过实际张数, 都要夹紧,
		// 否则会切出负数下标 / 越界 panic。
		n := r.FreeCount
		if n < 0 {
			n = 0
		}
		if n > total {
			n = total
		}
		d.Pics = d.Pics[:n]
	}
	return &service.DetailDTO{
		Photo: d, Playable: acc.Playable, NeedPay: acc.NeedPay, NeedVip: acc.NeedVip,
		Enough: acc.Enough, Reason: acc.Reason,
		PreviewCount: len(d.Pics), TotalCount: total,
	}, nil
}

func (s *sPhoto) Buy(ctx context.Context, userId, id int64) (float64, float64, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return 0, 0, err
	}
	if r.IsVip == 1 {
		return 0, 0, gerror.New("该图集为会员专享, 请开通会员")
	}
	if r.Price <= 0 {
		return 0, 0, gerror.New("该图集免费, 无需购买")
	}
	if err := paywall.Buy(ctx, userId, paywall.MediaPhoto, r.Id, r.Title, r.Price); err != nil {
		return 0, 0, err
	}
	_, _ = g.Model("photo_album").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"buy_count": &gdb.Counter{Field: "buy_count", Value: 1}}).Update()
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return r.Price, 0, nil
	}
	return r.Price, bal.Float64(), nil
}

// ---------------- 后台 ----------------

func (s *sPhoto) List(ctx context.Context, f service.ListFilter) ([]*service.PhotoDTO, int, error) {
	return s.query(ctx, f)
}

func normalize(in *service.SaveInput) error {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return gerror.New("标题不能为空")
	}
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	if in.Price < 0 {
		return gerror.New("价格不能为负")
	}
	if in.FreeCount < 0 {
		in.FreeCount = 0
	}
	return nil
}

func (s *sPhoto) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if err := normalize(&in); err != nil {
		return 0, err
	}
	return g.Model("photo_album").Ctx(ctx).Data(g.Map{
		"site_id": phSiteId, "title": in.Title, "cover": in.Cover, "intro": in.Intro,
		"category": in.Category, "tags": encodeJSON(in.Tags), "is_vip": in.IsVip,
		"price": in.Price, "free_count": in.FreeCount,
		"pics": encodeJSON(in.Pics), "pic_count": len(in.Pics),
		"rank": in.Rank, "status": in.Status, "publish_id": in.PublishId,
	}).InsertAndGetId()
}

func (s *sPhoto) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	data := g.Map{
		"is_vip": in.IsVip, "price": in.Price, "free_count": in.FreeCount,
		"rank": in.Rank, "status": in.Status, "updated_at": gtime.Now(),
	}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Cover != "" {
		data["cover"] = in.Cover
	}
	if in.Intro != "" {
		data["intro"] = in.Intro
	}
	if in.Category != "" {
		data["category"] = in.Category
	}
	if in.Tags != nil {
		data["tags"] = encodeJSON(in.Tags)
	}
	if in.Pics != nil { // 图片是整体覆盖语义, 不传就保持原样
		data["pics"] = encodeJSON(in.Pics)
		data["pic_count"] = len(in.Pics)
	}
	_, err := g.Model("photo_album").Ctx(ctx).
		Where("site_id", phSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sPhoto) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	// 图片就在主表 jsonb 里, 没有从表要连带清理, 单条删除即可。
	_, err := g.Model("photo_album").Ctx(ctx).
		Where("site_id", phSiteId).Where("id", id).Delete()
	return err
}

func (s *sPhoto) Audit(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.New("状态非法")
	}
	res, err := g.Model("photo_album").Ctx(ctx).
		Where("site_id", phSiteId).Where("id", id).
		Data(g.Map{"status": status, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("图集不存在")
	}
	return nil
}
