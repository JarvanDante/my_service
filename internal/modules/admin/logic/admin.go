// Package logic 后台管理员业务实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/admin/domain"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

type sAdmin struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IAdmin { return &sAdmin{repo: repo} }

func (s *sAdmin) Login(ctx context.Context, in service.LoginInput) (*service.LoginDTO, error) {
	if in.Username == "" || in.Password == "" {
		return nil, gerror.New("账号或密码不能为空")
	}
	a, err := s.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, gerror.New("账号或密码错误")
	}
	if a.Status != 1 {
		return nil, gerror.New("账号已被禁用")
	}
	if gmd5.MustEncryptString(in.Password+a.Salt) != a.Password {
		return nil, gerror.New("账号或密码错误")
	}
	_ = s.repo.UpdateLoginInfo(ctx, a.Id, in.Ip)
	token, err := kit.IssueAdminToken(ctx, a.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{Token: token, Admin: toInfo(a)}, nil
}

func (s *sAdmin) Logout(ctx context.Context, adminId int64) error {
	return kit.RevokeAdminByAdminId(ctx, adminId)
}

func (s *sAdmin) Info(ctx context.Context, adminId int64) (*service.AdminInfoDTO, error) {
	a, err := s.repo.FindById(ctx, adminId)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, gerror.New("管理员不存在")
	}
	return toInfo(a), nil
}

func toInfo(a *entity.AdminUser) *service.AdminInfoDTO {
	return &service.AdminInfoDTO{Id: a.Id, Username: a.Username, Nickname: a.Nickname, RoleId: a.RoleId}
}
