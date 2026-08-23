// Package logic 排行/热搜业务(移植自 tianbi rank/hotsearch)。
// 排行: Mongo 聚合改 PG GROUP BY(user_collect 点赞), Redis 缓存 60s, Redis 异常降级直查。
package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/ranks/service"
)

const (
	rkSiteId   = 1
	rkTopN     = 50
	rkCacheTTL = 60 // 秒
	rkHotLimit = 10 // 前台热搜默认条数
)

// normalizeHotCategory 后台写入/前台筛选用的分类码。
// 空或「通用」= 全站共用; 中文别名落到 comic/cartoon/novel/short/video。
func normalizeHotCategory(raw string) (string, error) {
	switch strings.TrimSpace(raw) {
	case "", "通用", "all", "_common":
		return "", nil
	case "comic", "漫画":
		return "comic", nil
	case "cartoon", "动漫":
		return "cartoon", nil
	case "novel", "小说":
		return "novel", nil
	case "short", "短剧":
		return "short", nil
	case "video", "视频":
		return "video", nil
	case "planet", "星球", "社区", "community":
		return "planet", nil
	default:
		return "", gerror.New("分类非法")
	}
}

func scanHotKeywords(ctx context.Context, category string, onlyCat bool, limit int) ([]string, error) {
	m := g.Model("hot_search").Ctx(ctx).Where("site_id", rkSiteId).Where("status", 1)
	if onlyCat {
		m = m.Where("category", category)
	}
	var list []*entity.HotSearch
	if err := m.OrderDesc("heat").OrderDesc("search_count").Limit(limit).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(list))
	for _, r := range list {
		if r.Keyword != "" {
			out = append(out, r.Keyword)
		}
	}
	return out, nil
}

type sRank struct{}

func New() service.IRank { return &sRank{} }

func cacheKey(mediaType int, period string) string {
	return fmt.Sprintf("rank:%d:%d:%s", rkSiteId, mediaType, period)
}

// Rank 点赞聚合排行 + Redis 缓存。
func (s *sRank) Rank(ctx context.Context, mediaType int, period string) ([]service.RankItem, error) {
	if period != "day" && period != "week" {
		period = "all"
	}
	key := cacheKey(mediaType, period)
	// 1. 读缓存(异常降级)
	if v, err := g.Redis().Get(ctx, key); err == nil && !v.IsNil() && v.String() != "" {
		var cached []service.RankItem
		if json.Unmarshal([]byte(v.String()), &cached) == nil {
			return cached, nil
		}
	}
	// 2. PG GROUP BY 聚合(tianbi Mongo $group 的等价物)
	m := g.Model("user_collect").Ctx(ctx).
		Where("site_id", rkSiteId).Where("op_type", 2). // 点赞
		Where("media_type", mediaType)
	switch period {
	case "day":
		m = m.Where("created_at >= ?", gtime.Now().StartOfDay())
	case "week":
		m = m.Where("created_at >= ?", gtime.Now().AddDate(0, 0, -7))
	}
	all, err := m.Fields("content_id, COUNT(*) AS score").
		Group("content_id").Order("score DESC, content_id DESC").Limit(rkTopN).All()
	if err != nil {
		return nil, err
	}
	out := make([]service.RankItem, 0, len(all))
	for i, rec := range all {
		out = append(out, service.RankItem{
			ContentId: rec["content_id"].Int64(), MediaType: mediaType,
			Score: rec["score"].Int64(), RankNo: i + 1,
		})
	}
	// 3. 写缓存(尽力而为)
	if b, e := json.Marshal(out); e == nil {
		_ = g.Redis().SetEX(ctx, key, string(b), rkCacheTTL)
	}
	return out, nil
}

// RefreshRank 清除全部排行缓存。
func (s *sRank) RefreshRank(ctx context.Context) error {
	for _, mt := range []int{1, 2} {
		for _, p := range []string{"day", "week", "all"} {
			_, _ = g.Redis().Del(ctx, cacheKey(mt, p))
		}
	}
	return nil
}

// HotKeywords 前台热搜词, 默认 10 条。
// category 非空: 先取该分类, 不够再用通用词补齐; 空: 全站混排(兼容旧端)。
func (s *sRank) HotKeywords(ctx context.Context, category string) ([]string, error) {
	cat, err := normalizeHotCategory(category)
	if err != nil {
		cat = ""
	}
	if cat == "" && strings.TrimSpace(category) == "" {
		return scanHotKeywords(ctx, "", false, rkHotLimit)
	}
	out, err := scanHotKeywords(ctx, cat, true, rkHotLimit)
	if err != nil {
		return nil, err
	}
	if len(out) >= rkHotLimit {
		return out[:rkHotLimit], nil
	}
	seen := make(map[string]struct{}, len(out))
	for _, w := range out {
		seen[w] = struct{}{}
	}
	extra, err := scanHotKeywords(ctx, "", true, rkHotLimit)
	if err != nil {
		return out, nil
	}
	for _, w := range extra {
		if _, ok := seen[w]; ok {
			continue
		}
		out = append(out, w)
		if len(out) >= rkHotLimit {
			break
		}
	}
	return out, nil
}

func (s *sRank) HotList(ctx context.Context, f service.HotFilter) ([]*service.HotDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("hot_search").Ctx(ctx).Where("site_id", rkSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Category != "" {
		if f.Category == "_common" {
			m = m.Where("category", "")
		} else if cat, e := normalizeHotCategory(f.Category); e == nil {
			m = m.Where("category", cat)
		}
	}
	if f.Keyword != "" {
		m = m.Where("keyword ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.HotSearch
	if err := m.Clone().OrderDesc("heat").OrderDesc("search_count").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.HotDTO, 0, len(list))
	for _, r := range list {
		updated := ""
		if r.UpdatedAt != nil {
			updated = r.UpdatedAt.String()
		}
		out = append(out, &service.HotDTO{
			Id: r.Id, Keyword: r.Keyword, Category: r.Category, Heat: r.Heat,
			SearchCount: r.SearchCount, Status: r.Status, UpdatedAt: updated,
		})
	}
	return out, total, nil
}

func (s *sRank) HotCreate(ctx context.Context, keyword, category string, heat, status int) (int64, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return 0, gerror.New("关键词不能为空")
	}
	cat, err := normalizeHotCategory(category)
	if err != nil {
		return 0, err
	}
	cnt, err := g.Model("hot_search").Ctx(ctx).
		Where("site_id", rkSiteId).Where("category", cat).Where("keyword", keyword).Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.New("该分类下关键词已存在")
	}
	if status != 0 && status != 1 {
		status = 1
	}
	id, err := g.Model("hot_search").Ctx(ctx).Data(g.Map{
		"site_id": rkSiteId, "keyword": keyword, "category": cat, "heat": heat, "status": status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sRank) HotUpdate(ctx context.Context, id int64, keyword, category string, heat, status int) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	cat, err := normalizeHotCategory(category)
	if err != nil {
		return err
	}
	nextKw := strings.TrimSpace(keyword)
	if nextKw != "" {
		cnt, e := g.Model("hot_search").Ctx(ctx).
			Where("site_id", rkSiteId).Where("category", cat).Where("keyword", nextKw).
			WhereNot("id", id).Count()
		if e != nil {
			return e
		}
		if cnt > 0 {
			return gerror.New("该分类下关键词已存在")
		}
	}
	data := g.Map{"heat": heat, "category": cat, "updated_at": gtime.Now()}
	if nextKw != "" {
		data["keyword"] = nextKw
	}
	if status == 0 || status == 1 {
		data["status"] = status
	}
	_, err = g.Model("hot_search").Ctx(ctx).
		Where("site_id", rkSiteId).Where("id", id).Data(data).Update()
	return err
}

func (s *sRank) HotDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("hot_search").Ctx(ctx).
		Where("site_id", rkSiteId).Where("id", id).Delete()
	return err
}
