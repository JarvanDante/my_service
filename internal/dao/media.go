package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

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

func (r *mediaRepo) MultipartCreate(ctx context.Context, m *entity.MediaMultipart) (int64, error) {
	res, err := g.Model("media_multipart").Ctx(ctx).Data(g.Map{
		"upload_id":       m.UploadId,
		"minio_upload_id": m.MinioUploadId,
		"bucket":          m.Bucket,
		"object_key":      m.ObjectKey,
		"purpose":         m.Purpose,
		"filename":        m.Filename,
		"content_type":    m.ContentType,
		"size":            m.Size,
		"part_size":       m.PartSize,
		"part_count":      m.PartCount,
		"status":          m.Status,
		"created_by":      m.CreatedBy,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *mediaRepo) MultipartFindByUploadId(ctx context.Context, uploadId string) (*entity.MediaMultipart, error) {
	var m *entity.MediaMultipart
	err := g.Model("media_multipart").Ctx(ctx).Where("upload_id", uploadId).Scan(&m)
	return m, err
}

func (r *mediaRepo) MultipartUpdateStatus(ctx context.Context, uploadId string, status int, mediaId int64) error {
	data := g.Map{
		"status":     status,
		"updated_at": gtime.Now(),
	}
	if mediaId > 0 {
		data["media_id"] = mediaId
	}
	_, err := g.Model("media_multipart").Ctx(ctx).Where("upload_id", uploadId).Data(data).Update()
	return err
}
