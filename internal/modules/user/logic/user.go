// Package logic 用户业务实现。依赖 domain.Repository 接口而非具体 dao。
package logic

import (
	"context"

	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

type sUser struct {
	repo domain.Repository
}

// New 注入仓储, 返回 service.IUser 实现。
func New(repo domain.Repository) service.IUser {
	return &sUser{repo: repo}
}

func (s *sUser) GetUser(ctx context.Context, id int64) (*service.UserDTO, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &service.UserDTO{ID: u.ID, Name: u.Name, Active: u.IsActive()}, nil
}

func (s *sUser) DisableUser(ctx context.Context, id int64) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err = u.Disable(); err != nil { // 调用领域规则
		return err
	}
	return s.repo.Save(ctx, u)
}
