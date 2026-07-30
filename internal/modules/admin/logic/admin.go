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
	"github.com/JarvanDante/my_service/internal/shared/rbac"
)

const superAdminCode = "superadmin"

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

func (s *sAdmin) ListRoles(ctx context.Context) ([]*service.RoleDTO, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.RoleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, &service.RoleDTO{
			Id: r.Id, Name: r.Name, Code: r.Code, Remark: r.Remark, Status: r.Status,
		})
	}
	return out, nil
}

func (s *sAdmin) ListPerms(ctx context.Context, roleCode string) ([]*service.PermDTO, error) {
	if roleCode == "" {
		return nil, gerror.New("角色码不能为空")
	}
	rows := rbac.PermsForRole(roleCode)
	out := make([]*service.PermDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, &service.PermDTO{Path: r[0], Method: r[1]})
	}
	return out, nil
}

func (s *sAdmin) AddPerm(ctx context.Context, roleCode, path, method string) error {
	if roleCode == "" || path == "" || method == "" {
		return gerror.New("角色码/路径/方法均不能为空")
	}
	if roleCode == superAdminCode {
		return gerror.New("超级管理员无需配置权限")
	}
	_, err := rbac.AddPolicy(roleCode, path, method)
	return err
}

func (s *sAdmin) RemovePerm(ctx context.Context, roleCode, path, method string) error {
	if roleCode == "" || path == "" || method == "" {
		return gerror.New("角色码/路径/方法均不能为空")
	}
	_, err := rbac.RemovePolicy(roleCode, path, method)
	return err
}

func toInfo(a *entity.AdminUser) *service.AdminInfoDTO {
	return &service.AdminInfoDTO{Id: a.Id, Username: a.Username, Nickname: a.Nickname, RoleId: a.RoleId}
}
