package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/comics/service"
)

type sModule struct {
	comics service.IComics
}

func NewModule() service.IModule { return &sModule{comics: New()} }

func decodeI64s(raw string) []int64 {
	out := []int64{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func encodeI64s(ids []int64) string {
	if ids == nil {
		ids = []int64{}
	}
	b, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func normalizeStyle(style int) int {
	if style < 1 || style > 8 {
		return entity.ComicsModuleStylePoster3x3
	}
	return style
}

func normalizeIcon(icon int) int {
	if icon < 1 || icon > 3 {
		return entity.ComicsModuleIconNew
	}
	return icon
}

func (s *sModule) defaultCatPosition(ctx context.Context) string {
	var row struct {
		Id int64 `orm:"id"`
	}
	_ = g.Model("comics_category").Ctx(ctx).
		Where("site_id", cmSiteId).Where("status", 1).
		OrderDesc("rank").OrderDesc("id").Limit(1).Scan(&row)
	if row.Id <= 0 {
		return ""
	}
	return "cat_" + strconv.FormatInt(row.Id, 10)
}

func (s *sModule) resolvePosition(ctx context.Context, pos string) string {
	pos = strings.TrimSpace(pos)
	if pos == "" || pos == entity.ComicsModulePosHome {
		return s.defaultCatPosition(ctx)
	}
	return pos
}

func normalizeSize(n, style int) int {
	if n <= 0 {
		if style == entity.ComicsModuleStylePoster3x3 {
			return 9
		}
		if style == entity.ComicsModuleStylePoster3x2 {
			return 6
		}
		if style == entity.ComicsModuleStylePosterRail || style == entity.ComicsModuleStyleWideRail {
			return 10
		}
		return 6
	}
	if n > 30 {
		return 30
	}
	return n
}

func (s *sModule) categoryNames(ctx context.Context, ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}
	var rows []struct {
		Id   int64  `orm:"id"`
		Name string `orm:"name"`
	}
	_ = g.Model("comics_category").Ctx(ctx).
		Where("site_id", cmSiteId).WhereIn("id", ids).Scan(&rows)
	byID := make(map[int64]string, len(rows))
	for _, r := range rows {
		byID[r.Id] = r.Name
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if name := byID[id]; name != "" {
			out = append(out, name)
		}
	}
	return out
}

func (s *sModule) tagNames(ctx context.Context, ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}
	var rows []struct {
		Id   int64  `orm:"id"`
		Name string `orm:"name"`
	}
	_ = g.Model("tag").Ctx(ctx).
		Where("site_id", cmSiteId).Where("content_type", 4).
		WhereIn("id", ids).Scan(&rows)
	byID := make(map[int64]string, len(rows))
	for _, r := range rows {
		byID[r.Id] = r.Name
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if name := byID[id]; name != "" {
			out = append(out, name)
		}
	}
	return out
}

func toModuleDTO(r *entity.ComicsModule, catNames, tagNames []string, filter string) *service.ModuleDTO {
	created, updated := "", ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	if r.UpdatedAt != nil {
		updated = r.UpdatedAt.String()
	}
	if catNames == nil {
		catNames = []string{}
	}
	if tagNames == nil {
		tagNames = []string{}
	}
	return &service.ModuleDTO{
		Id: r.Id, Name: r.Name, Position: r.Position, Style: r.Style, Icon: r.Icon,
		CategoryIds: decodeI64s(r.CategoryIds), CategoryNames: catNames,
		TagIds: decodeI64s(r.TagIds), TagNames: tagNames, Filter: filter,
		Size: r.Size, Rank: r.Rank, Status: r.Status,
		CreatedAt: created, UpdatedAt: updated,
	}
}

func (s *sModule) bindFilter(in *service.ModuleInput) moduleQuery {
	q := parseModuleFilter(in.Filter, in.CategoryIds, in.TagIds)
	in.CategoryIds = q.CatIDs
	in.TagIds = q.TagIDs
	in.Filter = encodeFilter(q)
	return q
}

func (s *sModule) queryOf(ctx context.Context, r *entity.ComicsModule) (moduleQuery, []string, []string) {
	q := parseModuleFilter(r.Filter, decodeI64s(r.CategoryIds), decodeI64s(r.TagIds))
	if id := parseCatPosition(r.Position); id > 0 && len(q.CatIDs) == 0 {
		if name, kind := s.categoryKind(ctx, id); kind == entity.ComicsCategoryKindNormal && name != "" {
			q.CatIDs = []int64{id}
		}
	}
	return q, s.categoryNames(ctx, q.CatIDs), s.tagNames(ctx, q.TagIDs)
}

func (s *sModule) categoryKind(ctx context.Context, id int64) (string, int) {
	if id <= 0 {
		return "", -1
	}
	var row struct {
		Name string `orm:"name"`
		Kind int    `orm:"kind"`
	}
	_ = g.Model("comics_category").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", id).Scan(&row)
	return row.Name, row.Kind
}

func (s *sModule) List(ctx context.Context, f service.ModuleFilter) ([]*service.ModuleDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("comics_module").Ctx(ctx).Where("site_id", cmSiteId)
	if name := strings.TrimSpace(f.Name); name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	if pos := strings.TrimSpace(f.Position); pos != "" {
		m = m.Where("position", pos)
	}
	if f.CategoryId > 0 {
		m = m.Where("category_ids @> ?::jsonb", encodeI64s([]int64{f.CategoryId}))
	}
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ComicsModule
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ModuleDTO, 0, len(list))
	for _, r := range list {
		q := parseModuleFilter(r.Filter, decodeI64s(r.CategoryIds), decodeI64s(r.TagIds))
		out = append(out, toModuleDTO(r, s.categoryNames(ctx, q.CatIDs), s.tagNames(ctx, q.TagIDs), encodeFilter(q)))
	}
	return out, total, nil
}

func (s *sModule) Create(ctx context.Context, in service.ModuleInput) (int64, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, gerror.New("模块名不能为空")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	style := normalizeStyle(in.Style)
	s.bindFilter(&in)
	pos := s.resolvePosition(ctx, in.Position)
	if pos == "" {
		return 0, gerror.New("请先在漫画分类里配置分类")
	}
	return g.Model("comics_module").Ctx(ctx).Data(g.Map{
		"site_id":      cmSiteId,
		"name":         name,
		"position":     pos,
		"style":        style,
		"icon":         normalizeIcon(in.Icon),
		"category_ids": encodeI64s(in.CategoryIds),
		"tag_ids":      encodeI64s(in.TagIds),
		"filter":       in.Filter,
		"size":         normalizeSize(in.Size, style),
		"rank":         in.Rank,
		"status":       in.Status,
	}).InsertAndGetId()
}

func (s *sModule) Update(ctx context.Context, in service.ModuleInput) error {
	if in.Id <= 0 {
		return gerror.New("模块ID非法")
	}
	style := normalizeStyle(in.Style)
	s.bindFilter(&in)
	pos := s.resolvePosition(ctx, in.Position)
	if pos == "" {
		return gerror.New("请先在漫画分类里配置分类")
	}
	data := g.Map{
		"position":     pos,
		"style":        style,
		"icon":         normalizeIcon(in.Icon),
		"category_ids": encodeI64s(in.CategoryIds),
		"tag_ids":      encodeI64s(in.TagIds),
		"filter":       in.Filter,
		"size":         normalizeSize(in.Size, style),
		"rank":         in.Rank,
		"updated_at":   gtime.Now(),
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		data["name"] = name
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("comics_module").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sModule) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("模块ID非法")
	}
	_, err := g.Model("comics_module").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", id).Delete()
	return err
}

func (s *sModule) FrontRepo(ctx context.Context, position string) ([]*service.ModuleFrontDTO, error) {
	pos := s.resolvePosition(ctx, position)
	var list []*entity.ComicsModule
	err := g.Model("comics_module").Ctx(ctx).
		Where("site_id", cmSiteId).Where("status", 1).Where("position", pos).
		OrderDesc("rank").OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, err
	}
	out := make([]*service.ModuleFrontDTO, 0, len(list))
	for _, r := range list {
		q, catNames, tagNames := s.queryOf(ctx, r)
		size := normalizeSize(r.Size, r.Style)
		items, _, err := s.comics.FrontList(ctx, 0, q.listFilter(catNames, tagNames, size, false))
		if err != nil {
			items = nil
		}
		out = append(out, &service.ModuleFrontDTO{
			Id: r.Id, Name: r.Name, Style: normalizeStyle(r.Style), Icon: normalizeIcon(r.Icon),
			Size: size, Tags: tagNames, Categories: catNames, Items: items,
		})
	}
	return out, nil
}

func (s *sModule) pickShuffle(ctx context.Context, q moduleQuery, catNames, tagNames []string, size int, exclude []int64) ([]*service.ComicsDTO, error) {
	f := q.listFilter(catNames, tagNames, size, true)
	f.ExcludeIds = exclude
	items, _, err := s.comics.FrontList(ctx, 0, f)
	if err != nil {
		return nil, err
	}
	if len(items) >= size {
		return items[:size], nil
	}
	picked := make([]int64, 0, len(items))
	for _, it := range items {
		picked = append(picked, it.Id)
	}
	moreFilter := q.listFilter(catNames, tagNames, size-len(items), true)
	moreFilter.ExcludeIds = picked
	more, _, err := s.comics.FrontList(ctx, 0, moreFilter)
	if err != nil {
		return items, nil
	}
	return append(items, more...), nil
}

func (s *sModule) FrontRefresh(ctx context.Context, id int64, exclude []int64) (*service.ModuleFrontDTO, error) {
	if id <= 0 {
		return nil, gerror.New("模块ID非法")
	}
	var r *entity.ComicsModule
	err := g.Model("comics_module").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", id).Where("status", 1).Scan(&r)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, gerror.New("模块不存在")
	}
	q, catNames, tagNames := s.queryOf(ctx, r)
	size := normalizeSize(r.Size, r.Style)
	items, err := s.pickShuffle(ctx, q, catNames, tagNames, size, exclude)
	if err != nil {
		return nil, err
	}
	return &service.ModuleFrontDTO{
		Id: r.Id, Name: r.Name, Style: normalizeStyle(r.Style), Icon: normalizeIcon(r.Icon),
		Size: size, Tags: tagNames, Categories: catNames, Items: items,
	}, nil
}
