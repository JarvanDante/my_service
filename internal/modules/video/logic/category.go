package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
)

const vdSiteId = 1

type sCategory struct{ table string }

func NewCategory() service.ICategory { return NewCategoryTable("video_category") }

func NewCategoryTable(table string) service.ICategory {
	if table == "" {
		table = "video_category"
	}
	return &sCategory{table: table}
}

func (s *sCategory) List(ctx context.Context, f service.CategoryFilter) ([]*service.CategoryDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model(s.table).Ctx(ctx).Where("site_id", vdSiteId)
	if f.Kind >= 0 {
		m = m.Where("kind", f.Kind)
	}
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.VideoCategory
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.CategoryDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.CategoryDTO{
			Id: r.Id, Name: r.Name, Kind: r.Kind, Rank: r.Rank, Status: r.Status, CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sCategory) Create(ctx context.Context, in service.CategoryInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("分类名不能为空")
	}
	if in.Kind < 0 || in.Kind > 3 {
		return 0, gerror.New("类型非法")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	cnt, err := g.Model(s.table).Ctx(ctx).
		Where("site_id", vdSiteId).Where("name", in.Name).Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.New("已存在同名分类")
	}
	return g.Model(s.table).Ctx(ctx).Data(g.Map{
		"site_id": vdSiteId, "name": in.Name, "kind": in.Kind,
		"rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sCategory) Update(ctx context.Context, in service.CategoryInput) error {
	if in.Id <= 0 {
		return gerror.New("分类ID非法")
	}
	if in.Name != "" {
		cnt, err := g.Model(s.table).Ctx(ctx).
			Where("site_id", vdSiteId).Where("name", in.Name).WhereNot("id", in.Id).Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return gerror.New("已存在同名分类")
		}
	}
	data := g.Map{"rank": in.Rank, "kind": in.Kind, "updated_at": gtime.Now()}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model(s.table).Ctx(ctx).
		Where("site_id", vdSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sCategory) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("分类ID非法")
	}
	_, err := g.Model(s.table).Ctx(ctx).
		Where("site_id", vdSiteId).Where("id", id).Delete()
	return err
}
