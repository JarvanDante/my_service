// Package logic 短剧业务：整部默认付费 + 单集覆盖，购买按单集走 shared/paywall。
package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/drama/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
	"github.com/JarvanDante/my_service/internal/shared/paas"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const (
	dmSiteId     = 1
	dmPreviewSec = 5
	// 互动编码: user_collect.media_type=7, comment.media_type=6, content_id 均为 drama.id
	dmCollectMedia = 7
	dmCommentMedia = 6
)

type sDrama struct{}

func New() service.IDrama { return &sDrama{} }

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

func fmtTime(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.String()
}

const sortPublishedAt = 4 // FrontList 把 sort=0 改写成按上架时间倒序

// parsePublishedAt 解析后台提交的上架时间，空值表示未填写。
func parsePublishedAt(raw string) (*gtime.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil, nil
	}
	t := gtime.NewFromStr(s)
	if t == nil || t.IsZero() {
		return nil, gerror.New("上架时间格式非法, 如: 2026-11-27 15:04:05")
	}
	return t, nil
}

// firstPublishedAt 首次上架且尚无上架时间时补当前时间；已有值或未上架则保留原值。
func firstPublishedAt(online bool, incoming, old *gtime.Time) *gtime.Time {
	if incoming != nil && !incoming.IsZero() {
		return incoming
	}
	if !online {
		return nil
	}
	if old != nil && !old.IsZero() {
		return nil
	}
	return gtime.Now()
}

func resolveCategory(in *service.SaveInput) string {
	if in.Categories != nil {
		return strings.Join(kit.MergeNames(in.Categories), ",")
	}
	return strings.Join(kit.NamesCSV(in.Category), ",")
}

func toDTO(r *entity.Drama) *service.DramaDTO {
	return &service.DramaDTO{
		Id: r.Id, Title: r.Title, Cover: r.Cover, Intro: r.Intro,
		Category: r.Category, Categories: kit.NamesCSV(r.Category), Tags: decodeTags(r.Tags),
		IsVip: r.IsVip, Price: r.Price, FreeEpisode: r.FreeEpisode, EpisodeCount: r.EpisodeCount,
		ViewCount: r.ViewCount, BuyCount: r.BuyCount, LikeCount: r.LikeCount, UpdateStatus: r.UpdateStatus,
		Rank: r.Rank, IsRecommend: r.IsRecommend, Status: r.Status, MediaCode: r.MediaCode,
		PublishedAt: fmtTime(r.PublishedAt), CreatedAt: fmtTime(r.CreatedAt),
	}
}

func toEpisodeDTO(ctx context.Context, r *entity.DramaEpisode) *service.EpisodeDTO {
	return &service.EpisodeDTO{
		Id: r.Id, DramaId: r.DramaId, Seq: r.Seq, Title: r.Title, MediaCode: r.MediaCode,
		Duration: r.Duration, Cover: paas.CoverURL(ctx, r.MediaCode), PayMode: r.PayMode, Price: r.Price,
		Status: r.Status, MediaReady: r.MediaReady, CreatedAt: fmtTime(r.CreatedAt),
	}
}

// hydrateCover 本站未自定义封面时用媒资整部封面签名地址。
func hydrateCover(ctx context.Context, list []*service.DramaDTO) {
	for _, d := range list {
		if d == nil || d.MediaCode == "" || strings.Contains(d.Cover, "/my-storage/") {
			continue
		}
		if u := paas.CoverURL(ctx, d.MediaCode); u != "" {
			d.Cover = u
		}
	}
}

// episodeRule 单集生效的付费规则，返回值语义同 paywall.Check 的 isVip/price。
func episodeRule(d *entity.Drama, e *entity.DramaEpisode) (bool, float64) {
	switch e.PayMode {
	case entity.DramaPayFree:
		return false, 0
	case entity.DramaPayCoin:
		if e.Price > 0 {
			return false, e.Price
		}
		return false, d.Price
	case entity.DramaPayVip:
		return true, 0
	}
	if e.Seq <= d.FreeEpisode {
		return false, 0
	}
	if d.IsVip == 1 {
		return true, 0
	}
	return false, d.Price
}

func (s *sDrama) query(ctx context.Context, f service.ListFilter) ([]*service.DramaDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	if f.Size > 100 {
		f.Size = 100
	}
	base := g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId)
	if f.Status >= 0 {
		base = base.Where("status", f.Status)
	}
	if cates := kit.MergeNames(kit.NamesCSV(f.Category), f.Categories); len(cates) > 0 {
		base = base.Where(kit.CategoryOverlapWhere(), strings.Join(cates, ","))
	}
	if tagNames := kit.MergeNames(kit.NamesCSV(f.Tag), f.Tags); len(tagNames) > 0 {
		ors := make([]string, 0, len(tagNames))
		args := make([]any, 0, len(tagNames))
		for _, t := range tagNames {
			ors = append(ors, "tags @> ?::jsonb")
			args = append(args, encodeJSON([]string{t}))
		}
		base = base.Where("("+strings.Join(ors, " OR ")+")", args...)
	}
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		base = base.Where("title ILIKE ?", "%"+kw+"%")
	}
	if f.OnlyRecommend {
		base = base.Where("is_recommend", 1)
	}
	if f.UpdateStatus == 1 || f.UpdateStatus == 2 {
		base = base.Where("update_status", f.UpdateStatus)
	}
	if len(f.Ids) > 0 {
		base = base.WhereIn("id", f.Ids)
	}
	if len(f.ExcludeIds) > 0 {
		base = base.WhereNotIn("id", f.ExcludeIds)
	}
	switch f.PayType {
	case 1:
		base = base.Where("is_vip", 1)
	case 2:
		base = base.Where("is_vip", 0).Where("price > ?", 0)
	case 3:
		base = base.Where("is_vip", 0).Where("price <= ?", 0)
	}
	if f.Sort == 1 {
		base = base.Where("view_count > 0")
	}
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	m := base.Clone()
	var list []*entity.Drama
	if f.Shuffle {
		if err := m.OrderRandom().Limit(f.Size).Scan(&list); err != nil {
			return nil, 0, err
		}
	} else {
		switch f.Sort {
		case 1:
			m = m.OrderDesc("view_count")
		case 2:
			m = m.OrderDesc("updated_at")
		case 3:
			m = m.OrderDesc("like_count")
		case sortPublishedAt:
			m = m.Order("published_at DESC NULLS LAST")
		default:
			m = m.OrderDesc("rank")
		}
		if err := m.OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
			return nil, 0, err
		}
	}
	out := make([]*service.DramaDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	hydrateCover(ctx, out)
	return out, total, nil
}

func (s *sDrama) find(ctx context.Context, id int64, onlineOnly bool) (*entity.Drama, error) {
	m := g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("id", id)
	if onlineOnly {
		m = m.Where("status", entity.ContentStatusOnline)
	}
	var r *entity.Drama
	if err := m.Scan(&r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, gerror.New("短剧不存在或已下架")
	}
	return r, nil
}

// findEpisode 前台取可播剧集及其上架短剧。
func (s *sDrama) findEpisode(ctx context.Context, episodeId int64) (*entity.DramaEpisode, *entity.Drama, error) {
	var e *entity.DramaEpisode
	if err := g.Model("drama_episode").Ctx(ctx).
		Where("site_id", dmSiteId).Where("id", episodeId).Where("status", 1).Where("media_ready", 1).
		Scan(&e); err != nil {
		return nil, nil, err
	}
	if e == nil {
		return nil, nil, gerror.New("剧集不存在或已下架")
	}
	d, err := s.find(ctx, e.DramaId, true)
	if err != nil {
		return nil, nil, err
	}
	return e, d, nil
}

func (s *sDrama) FrontList(ctx context.Context, f service.ListFilter) ([]*service.DramaDTO, int, error) {
	f.Status = entity.ContentStatusOnline
	if f.Sort == 0 {
		f.Sort = sortPublishedAt
	}
	return s.query(ctx, f)
}

func (s *sDrama) Detail(ctx context.Context, id int64) (*service.DramaDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return nil, err
	}
	_, _ = g.Model("drama").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"view_count": &gdb.Counter{Field: "view_count", Value: 1}}).Update()
	d := toDTO(r)
	hydrateCover(ctx, []*service.DramaDTO{d})
	d.ViewCount++
	return d, nil
}

func (s *sDrama) Interact(ctx context.Context, userId, id int64) (*service.InteractDTO, error) {
	out := &service.InteractDTO{}
	n, err := g.Model("comment").Ctx(ctx).Where("site_id", dmSiteId).Where("media_type", dmCommentMedia).
		Where("content_id", id).Where("status", 1).Count()
	if err != nil {
		return nil, err
	}
	out.CommentCount = n
	if userId <= 0 {
		return out, nil
	}
	ops, err := g.Model("user_collect").Ctx(ctx).Where("site_id", dmSiteId).Where("user_id", userId).
		Where("media_type", dmCollectMedia).Where("content_id", id).WhereIn("op_type", g.Slice{1, 2}).
		Fields("op_type").Array()
	if err != nil {
		return nil, err
	}
	for _, v := range ops {
		switch v.Int() {
		case 1:
			out.IsCollected = true
		case 2:
			out.IsLiked = true
		}
	}
	return out, nil
}

func (s *sDrama) Episodes(ctx context.Context, userId, id int64, desc bool) (string, []*service.EpisodeDTO, error) {
	d, err := s.find(ctx, id, true)
	if err != nil {
		return "", nil, err
	}
	m := g.Model("drama_episode").Ctx(ctx).
		Where("site_id", dmSiteId).Where("drama_id", id).Where("status", 1).Where("media_ready", 1)
	if desc {
		m = m.OrderDesc("seq")
	} else {
		m = m.OrderAsc("seq")
	}
	var list []*entity.DramaEpisode
	if err := m.Scan(&list); err != nil {
		return "", nil, err
	}
	out, err := s.applyAccess(ctx, userId, d, list)
	if err != nil {
		return "", nil, err
	}
	return d.Title, out, nil
}

// applyAccess 批量判定解锁态，规则与 paywall.Check 一致，只把会员和已购查询合并成各一次。
func (s *sDrama) applyAccess(ctx context.Context, userId int64, d *entity.Drama, list []*entity.DramaEpisode) ([]*service.EpisodeDTO, error) {
	paidIds := make([]int64, 0, len(list))
	hasVip := false
	for _, e := range list {
		isVip, price := episodeRule(d, e)
		if isVip {
			hasVip = true
		}
		if isVip || price > 0 {
			paidIds = append(paidIds, e.Id)
		}
	}
	vipActive := false
	if userId > 0 && hasVip {
		ok, err := paywall.IsVipActive(ctx, userId)
		if err != nil {
			return nil, err
		}
		vipActive = ok
	}
	bought, err := paywall.PurchasedSet(ctx, userId, paywall.MediaDrama, paidIds)
	if err != nil {
		return nil, err
	}
	out := make([]*service.EpisodeDTO, 0, len(list))
	for _, e := range list {
		dto := toEpisodeDTO(ctx, e)
		isVip, price := episodeRule(d, e)
		dto.IsVip, dto.Price = isVip, price
		switch {
		case (isVip || price > 0) && bought[e.Id]:
			// 已购优先：买过的集无论之后改成 VIP 还是改价都可看全片。
			dto.IsBuy, dto.Playable = true, true
		case isVip:
			dto.Playable, dto.NeedVip = vipActive, !vipActive
		case price > 0:
			dto.IsBuy = bought[e.Id]
			dto.Playable, dto.NeedPay = dto.IsBuy, !dto.IsBuy
		default:
			dto.IsFree, dto.Playable, dto.Price = true, true, 0
		}
		out = append(out, dto)
	}
	return out, nil
}

func (s *sDrama) Play(ctx context.Context, userId, episodeId int64) (*service.PlayDTO, error) {
	e, d, err := s.findEpisode(ctx, episodeId)
	if err != nil {
		return nil, err
	}
	isVip, price := episodeRule(d, e)
	acc := &paywall.Access{Price: price}
	// 已购优先：VIP 集先查购买记录，金币集由 paywall.Check 自身先查已购。
	bought := false
	if isVip && userId > 0 {
		if bought, err = paywall.Purchased(ctx, userId, paywall.MediaDrama, e.Id); err != nil {
			return nil, err
		}
	}
	if bought {
		acc.Playable, acc.IsBuy = true, true
	} else if acc, err = paywall.Check(ctx, userId, paywall.MediaDrama, e.Id, isVip, price); err != nil {
		return nil, err
	}
	preview := 0
	if !acc.Playable {
		preview = dmPreviewSec
	}
	list, err := paas.DramaPlaylist(ctx, d.MediaCode, e.MediaCode)
	if err != nil {
		g.Log().Warningf(ctx, "短剧播放授权失败 episode=%d code=%s: %v", e.Id, e.MediaCode, err)
		return nil, gerror.New("播放地址获取失败")
	}
	playUrl := paas.PlaylistURL(ctx, e.MediaCode, "/"+list, preview)
	if playUrl == "/"+list {
		if playUrl, err = paas.DramaPlayToken(ctx, e.MediaCode, preview); err != nil {
			g.Log().Warningf(ctx, "短剧播放签发失败 episode=%d code=%s: %v", e.Id, e.MediaCode, err)
			return nil, gerror.New("播放地址获取失败")
		}
	}
	dto := toEpisodeDTO(ctx, e)
	dto.IsVip, dto.Price, dto.IsFree = isVip, price, !isVip && price <= 0
	dto.Playable, dto.IsBuy, dto.NeedPay, dto.NeedVip = acc.Playable, acc.IsBuy, acc.NeedPay, acc.NeedVip
	out := &service.PlayDTO{
		Episode: dto, PlayUrl: playUrl,
		PreviewSec: preview, Enough: acc.Enough, Reason: acc.Reason,
	}
	near := g.Model("drama_episode").Ctx(ctx).
		Where("site_id", dmSiteId).Where("drama_id", e.DramaId).Where("status", 1).Where("media_ready", 1)
	if v, _ := near.Clone().Where("seq < ?", e.Seq).OrderDesc("seq").Fields("id").Value(); v != nil {
		out.PrevId = v.Int64()
	}
	if v, _ := near.Clone().Where("seq > ?", e.Seq).OrderAsc("seq").Fields("id").Value(); v != nil {
		out.NextId = v.Int64()
	}
	return out, nil
}

func (s *sDrama) Buy(ctx context.Context, userId, episodeId int64) (float64, float64, error) {
	e, d, err := s.findEpisode(ctx, episodeId)
	if err != nil {
		return 0, 0, err
	}
	isVip, price := episodeRule(d, e)
	if isVip {
		return 0, 0, gerror.New("该集为会员专享, 请开通会员")
	}
	if price <= 0 {
		return 0, 0, gerror.New("该集免费, 无需购买")
	}
	title := d.Title + " 第" + strconv.Itoa(e.Seq) + "集"
	if err := paywall.Buy(ctx, userId, paywall.MediaDrama, e.Id, title, price); err != nil {
		return 0, 0, err
	}
	_, _ = g.Model("drama").Ctx(ctx).Where("id", d.Id).
		Data(g.Map{"buy_count": &gdb.Counter{Field: "buy_count", Value: 1}}).Update()
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil || bal == nil {
		return price, 0, nil
	}
	return price, bal.Float64(), nil
}
