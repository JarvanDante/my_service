// Package domain 用户领域层。
package domain

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	FindById(ctx context.Context, id int64) (*entity.Users, error)
	FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Users, error)
	Create(ctx context.Context, data g.Map) (int64, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
	UpdatePhone(ctx context.Context, id int64, phone string) error
	Disable(ctx context.Context, id int64, reason string) error
}
