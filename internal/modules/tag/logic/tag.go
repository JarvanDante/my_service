// Package logic 标签业务(移植自 tianbi tagser)。
// tianbi 按 Type 分派到 mediatag/comicstag/actressphototag/noveltag/postdiscusstag 多个集合读取,
// PG 统一为单表 tag + content_type 列, 一处读写。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/tag/service"
)

const (
	tagSiteId  = 1  // 单站点样板
	repoMaxTag = 50 // 前台单次最多返回(与 tianbi 上限一致)
)

type sTag struct{}

func New() service.ITag { return &sTag{} }

// Repo 前台按内容类型取启用标签, rank desc(移植自 tianbi TagDetail)。
func (s *sTag) Repo(ctx context.Context, contentType, page, size int) ([]service.RepoItem, error) {
	if contentType <= 0 {
		return nil, gerror.New("内容类型非法")
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > repoMaxTag {
		size = repoMaxTag
	}
	var list []*entity.Tag
	err := g.Model("tag").Ctx(ctx).
		Where("site_id", tagSiteId).
		Where("content_type", contentType).
		Where("status", 1).
		OrderDesc("rank").OrderDesc("id").
		Page(page, size).Scan(&list)
	if err != nil {
		return nil, err
	}
	out := make([]service.RepoItem, 0, len(list))
	for _, r := range list {
		out = append(out, service.RepoItem{Id: r.Id, Name: r.Name})
	}
	return out, nil
}

func (s *sTag) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("tag").Ctx(ctx).Where("site_id", tagSiteId)
	if f.ContentType > 0 {
		m = m.Where("content_type", f.ContentType)
	}
	if f.Status >= 0 { // -1 = 全部; 0/1 = 精确过滤
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Tag
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.ItemDTO{
			Id: r.Id, ContentType: r.ContentType, Name: r.Name,
			Rank: r.Rank, Status: r.Status, CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sTag) Create(ctx context.Context, in service.CreateInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("标签名不能为空")
	}
	if in.ContentType <= 0 {
		return 0, gerror.New("内容类型非法")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	// 同类型下名称唯一
	cnt, err := g.Model("tag").Ctx(ctx).
		Where("site_id", tagSiteId).Where("content_type", in.ContentType).
		Where("name", in.Name).Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.New("该类型下已存在同名标签")
	}
	id, err := g.Model("tag").Ctx(ctx).Data(g.Map{
		"site_id": tagSiteId, "content_type": in.ContentType,
		"name": in.Name, "rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sTag) Update(ctx context.Context, in service.UpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("标签ID非法")
	}
	data := g.Map{"rank": in.Rank, "updated_at": gtime.Now()}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("tag").Ctx(ctx).
		Where("site_id", tagSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sTag) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("标签ID非法")
	}
	_, err := g.Model("tag").Ctx(ctx).
		Where("site_id", tagSiteId).Where("id", id).Delete()
	return err
}
