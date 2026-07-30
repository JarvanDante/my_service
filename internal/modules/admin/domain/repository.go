// Package domain 后台管理员领域层。
package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*entity.AdminUser, error)
	FindById(ctx context.Context, id int64) (*entity.AdminUser, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
}
