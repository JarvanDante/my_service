package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	videodomain "github.com/JarvanDante/my_service/internal/modules/video/domain"
)

type videoRepo struct{}

func NewVideoRepo() videodomain.Repository { return &videoRepo{} }

func (r *videoRepo) List(ctx context.Context, f videodomain.ListFilter, page, size int) ([]*entity.Video, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("video").Ctx(ctx)
	if f.Keyword != "" {
		m = m.WhereLike("title", "%"+f.Keyword+"%")
	}
	if f.MediaCode != "" {
		m = m.Where("media_code", f.MediaCode)
	}
	if f.Category != "" {
		m = m.Where("string_to_array(replace(category, '，', ','), ',') @> ARRAY[?]::text[]", f.Category)
	}
	m = m.Where("kind", f.Kind)
	if f.Status != 9 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	q := m.Clone()
	// 排序在 SQL 里做, 保证与分页一致。default 分支即原有的"综合"顺序, 后台不传 Sort 时行为不变。
	switch f.Sort {
	case 1: // 最新: 自增主键即时间序, 比按 created_at 排省一个索引
		q = q.OrderDesc("id")
	case 2: // 时长
		q = q.OrderDesc("duration").OrderDesc("id")
	default:
		q = q.OrderDesc("sort").OrderDesc("id")
	}
	var list []*entity.Video
	err = q.Page(page, size).Scan(&list)
	return list, total, err
}

func (r *videoRepo) Find(ctx context.Context, id int64) (*entity.Video, error) {
	var v *entity.Video
	err := g.Model("video").Ctx(ctx).Where("id", id).Scan(&v)
	return v, err
}

func (r *videoRepo) FindByMediaCode(ctx context.Context, code string, kind int) (*entity.Video, error) {
	var v *entity.Video
	err := g.Model("video").Ctx(ctx).Where("media_code", code).Where("kind", kind).Scan(&v)
	return v, err
}

func (r *videoRepo) Create(ctx context.Context, v *entity.Video) (int64, error) {
	res, err := g.Model("video").Ctx(ctx).Data(g.Map{
		"title": v.Title, "description": v.Description,
		"cover_url": v.CoverUrl, "cover_key": v.CoverKey, "cover_media_id": v.CoverMediaId,
		"source_url": v.SourceUrl, "source_key": v.SourceKey, "source_media_id": v.SourceMediaId,
		"media_code": v.MediaCode, "kind": v.Kind, "category": v.Category, "tags": v.Tags,
		"duration": v.Duration, "sort": v.Sort, "status": v.Status, "created_by": v.CreatedBy,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *videoRepo) Update(ctx context.Context, v *entity.Video) error {
	_, err := g.Model("video").Ctx(ctx).Where("id", v.Id).Data(g.Map{
		"title": v.Title, "description": v.Description,
		"cover_url": v.CoverUrl, "cover_key": v.CoverKey, "cover_media_id": v.CoverMediaId,
		"source_url": v.SourceUrl, "source_key": v.SourceKey, "source_media_id": v.SourceMediaId,
		"media_code": v.MediaCode, "kind": v.Kind, "category": v.Category, "tags": v.Tags,
		"duration": v.Duration, "sort": v.Sort, "status": v.Status,
		"updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *videoRepo) Delete(ctx context.Context, id int64) error {
	_, err := g.Model("video").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *videoRepo) SetStatus(ctx context.Context, id int64, status int) error {
	_, err := g.Model("video").Ctx(ctx).Where("id", id).Data(g.Map{
		"status": status, "updated_at": gtime.Now(),
	}).Update()
	return err
}
