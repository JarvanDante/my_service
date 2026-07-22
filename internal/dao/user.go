// Package dao 数据访问层: 仓储接口的 PostgreSQL 实现。g.Model/gdb 仅出现在此包。
package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
)

type userRepo struct{}

// NewUserRepo 返回 user 领域的仓储实现。
func NewUserRepo() domain.Repository {
	return &userRepo{}
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var po entity.User
	if err := g.Model("users").Ctx(ctx).WherePri(id).Scan(&po); err != nil {
		return nil, err
	}
	return &domain.User{
		ID:     po.Id,
		Name:   po.Name,
		Email:  po.Email,
		Status: domain.Status(po.Status),
	}, nil
}

func (r *userRepo) Save(ctx context.Context, u *domain.User) error {
	_, err := g.Model("users").Ctx(ctx).WherePri(u.ID).Data(g.Map{
		"name":   u.Name,
		"email":  u.Email,
		"status": u.Status,
	}).Update()
	return err
}
