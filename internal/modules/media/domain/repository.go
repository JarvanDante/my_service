package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	Create(ctx context.Context, m *entity.MediaObject) (int64, error)

	MultipartCreate(ctx context.Context, m *entity.MediaMultipart) (int64, error)
	MultipartFindByUploadId(ctx context.Context, uploadId string) (*entity.MediaMultipart, error)
	MultipartUpdateStatus(ctx context.Context, uploadId string, status int, mediaId int64) error
}
