package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type ListFilter struct {
	Keyword string
	Status  int // 9=全部
}

type Repository interface {
	List(ctx context.Context, f ListFilter, page, size int) ([]*entity.Video, int, error)
	Find(ctx context.Context, id int64) (*entity.Video, error)
	Create(ctx context.Context, v *entity.Video) (int64, error)
	Update(ctx context.Context, v *entity.Video) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status int) error
}
