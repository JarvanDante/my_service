package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	Create(ctx context.Context, m *entity.MediaObject) (int64, error)
}
