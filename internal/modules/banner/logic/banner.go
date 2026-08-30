package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/banner/service"
)

const bannerSiteId = 1

var bannerPositions = map[string]string{
	entity.BannerPosComics:  "漫画",
	entity.BannerPosCartoon: "动漫",
	entity.BannerPosVideo:   "视频",
	entity.BannerPosNovel:   "小说",
}

type sBanner struct{}

func New() service.IBanner { return &sBanner{} }

func normalizePosition(raw string) (string, error) {
	pos := strings.ToLower(strings.TrimSpace(raw))
	if _, ok := bannerPositions[pos]; !ok {
		return "", gerror.New("位置非法")
	}
	return pos, nil
}

func toDTO(r *entity.Banner) *service.ItemDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.ItemDTO{
		Id: r.Id, Position: r.Position, Title: r.Title, CoverUrl: r.CoverUrl,
		Link: r.Link, Rank: r.Rank, Status: r.Status, CreatedAt: created,
	}
}

func (s *sBanner) FrontList(ctx context.Context, position string) ([]*service.ItemDTO, error) {
	pos, err := normalizePosition(position)
	if err != nil {
		return nil, err
	}
	var list []*entity.Banner
	if err := g.Model("banner").Ctx(ctx).
		Where("site_id", bannerSiteId).Where("position", pos).Where("status", 1).
		OrderDesc("rank").OrderDesc("id").Limit(20).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

func (s *sBanner) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("banner").Ctx(ctx).Where("site_id", bannerSiteId)
	if f.Position != "" {
		pos, err := normalizePosition(f.Position)
		if err != nil {
			return nil, 0, err
		}
		m = m.Where("position", pos)
	}
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Banner
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sBanner) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	pos, err := normalizePosition(in.Position)
	if err != nil {
		return 0, err
	}
	cover := strings.TrimSpace(in.CoverUrl)
	if cover == "" {
		return 0, gerror.New("请上传轮播图")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	return g.Model("banner").Ctx(ctx).Data(g.Map{
		"site_id": bannerSiteId, "position": pos, "title": strings.TrimSpace(in.Title),
		"cover_url": cover, "link": strings.TrimSpace(in.Link),
		"rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sBanner) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	pos, err := normalizePosition(in.Position)
	if err != nil {
		return err
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	data := g.Map{
		"position": pos, "title": strings.TrimSpace(in.Title),
		"link": strings.TrimSpace(in.Link), "rank": in.Rank,
		"status": in.Status, "updated_at": gtime.Now(),
	}
	if cover := strings.TrimSpace(in.CoverUrl); cover != "" {
		data["cover_url"] = cover
	}
	_, err = g.Model("banner").Ctx(ctx).
		Where("site_id", bannerSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sBanner) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("banner").Ctx(ctx).
		Where("site_id", bannerSiteId).Where("id", id).Delete()
	return err
}
