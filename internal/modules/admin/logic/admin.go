// Package logic 后台管理员业务实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/admin/domain"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

const superAdminCode = "superadmin"

type sAdmin struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IAdmin { return &sAdmin{repo: repo} }

// ---------- 认证 ----------

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

// ---------- 角色权限(Casbin) ----------

func (s *sAdmin) ListRoles(ctx context.Context) ([]*service.RoleDTO, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.RoleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, &service.RoleDTO{
			Id: r.Id, Name: r.Name, Code: r.Code, Remark: r.Remark, Status: r.Status, Permissions: r.Permissions,
		})
	}
	return out, nil
}

// ---------- 角色管理 ----------

func (s *sAdmin) CreateRole(ctx context.Context, in service.RoleCreateInput) (int64, error) {
	if in.Name == "" || in.Code == "" {
		return 0, gerror.New("角色名与角色码必填")
	}
	if in.Code == superAdminCode {
		return 0, gerror.New("superadmin 为系统保留角色码")
	}
	exist, err := s.repo.FindRoleByCode(ctx, in.Code)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("角色码已存在")
	}
	return s.repo.CreateRole(ctx, &entity.AdminRole{Name: in.Name, Code: in.Code, Remark: in.Remark, Permissions: in.Permissions})
}

func (s *sAdmin) UpdateRole(ctx context.Context, in service.RoleUpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("角色ID无效")
	}
	role, err := s.repo.FindRoleById(ctx, in.Id)
	if err != nil {
		return err
	}
	if role == nil {
		return gerror.New("角色不存在")
	}
	if role.Code == superAdminCode {
		return gerror.New("超级管理员角色不可修改")
	}
	return s.repo.UpdateRole(ctx, in.Id, in.Name, in.Remark, in.Status, in.Permissions)
}

func (s *sAdmin) DeleteRole(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("角色ID无效")
	}
	role, err := s.repo.FindRoleById(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return gerror.New("角色不存在")
	}
	if role.Code == superAdminCode {
		return gerror.New("超级管理员角色不可删除")
	}
	cnt, err := s.repo.CountAdminsByRoleId(ctx, id)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return gerror.New("该角色下仍有管理员, 不能删除")
	}
	return s.repo.DeleteRole(ctx, id)
}

// ---------- 管理员账号管理 ----------

func (s *sAdmin) ListAdmins(ctx context.Context, page, size int) (*service.AdminListDTO, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := s.repo.ListAdmins(ctx, page, size)
	if err != nil {
		return nil, err
	}
	// 角色 id → name 映射
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleName := make(map[int64]string, len(roles))
	for _, r := range roles {
		roleName[r.Id] = r.Name
	}
	items := make([]*service.AdminItemDTO, 0, len(list))
	for _, a := range list {
		last := ""
		if a.LastLoginAt != nil {
			last = a.LastLoginAt.String()
		}
		items = append(items, &service.AdminItemDTO{
			Id: a.Id, Username: a.Username, Nickname: a.Nickname,
			RoleId: a.RoleId, RoleName: roleName[a.RoleId], Status: a.Status, LastLoginAt: last,
		})
	}
	return &service.AdminListDTO{List: items, Total: total, Page: page, Size: size}, nil
}

func (s *sAdmin) CreateAdmin(ctx context.Context, in service.AdminCreateInput) (int64, error) {
	if in.Username == "" || in.Password == "" {
		return 0, gerror.New("账号与密码必填")
	}
	if in.RoleId <= 0 {
		return 0, gerror.New("请指定角色")
	}
	exist, err := s.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("账号已存在")
	}
	role, err := s.repo.FindRoleById(ctx, in.RoleId)
	if err != nil {
		return 0, err
	}
	if role == nil {
		return 0, gerror.New("角色不存在")
	}
	salt := grand.S(8)
	pwd := gmd5.MustEncryptString(in.Password + salt)
	return s.repo.CreateAdmin(ctx, &entity.AdminUser{
		Username: in.Username, Password: pwd, Salt: salt, Nickname: in.Nickname, RoleId: in.RoleId,
	})
}

func (s *sAdmin) UpdateAdmin(ctx context.Context, in service.AdminUpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("管理员ID无效")
	}
	a, err := s.repo.FindById(ctx, in.Id)
	if err != nil {
		return err
	}
	if a == nil {
		return gerror.New("管理员不存在")
	}
	if in.RoleId <= 0 {
		return gerror.New("请指定角色")
	}
	role, err := s.repo.FindRoleById(ctx, in.RoleId)
	if err != nil {
		return err
	}
	if role == nil {
		return gerror.New("角色不存在")
	}
	var pwd, salt string
	if in.Password != "" {
		salt = grand.S(8)
		pwd = gmd5.MustEncryptString(in.Password + salt)
	}
	return s.repo.UpdateAdmin(ctx, in.Id, in.Nickname, in.RoleId, in.Status, pwd, salt)
}

func (s *sAdmin) DeleteAdmin(ctx context.Context, id, operatorId int64) error {
	if id <= 0 {
		return gerror.New("管理员ID无效")
	}
	if id == operatorId {
		return gerror.New("不能删除当前登录账号")
	}
	a, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return gerror.New("管理员不存在")
	}
	return s.repo.DeleteAdmin(ctx, id)
}

func toInfo(a *entity.AdminUser) *service.AdminInfoDTO {
	return &service.AdminInfoDTO{Id: a.Id, Username: a.Username, Nickname: a.Nickname, RoleId: a.RoleId}
}
