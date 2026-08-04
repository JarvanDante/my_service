package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	mediadomain "github.com/JarvanDante/my_service/internal/modules/media/domain"
)

type mediaRepo struct{}

func NewMediaRepo() mediadomain.Repository { return &mediaRepo{} }

func (r *mediaRepo) Create(ctx context.Context, m *entity.MediaObject) (int64, error) {
	res, err := g.Model("media_object").Ctx(ctx).Data(g.Map{
		"bucket":       m.Bucket,
		"object_key":   m.ObjectKey,
		"purpose":      m.Purpose,
		"content_type": m.ContentType,
		"size":         m.Size,
		"etag":         m.Etag,
		"created_by":   m.CreatedBy,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
