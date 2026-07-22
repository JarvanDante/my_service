package domain

import "context"

// Repository 仓储接口。实现放在 internal/dao, 依赖倒置。
type Repository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	Save(ctx context.Context, u *User) error
}
