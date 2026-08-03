package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	admindomain "github.com/JarvanDante/my_service/internal/modules/admin/domain"
)

type adminRepo struct{}

func NewAdminRepo() admindomain.Repository { return &adminRepo{} }

// ---------- 认证 ----------

func (r *adminRepo) FindByUsername(ctx context.Context, username string) (*entity.AdminUser, error) {
	var a *entity.AdminUser
	err := g.Model("admin_user").Ctx(ctx).Where("username", username).Scan(&a)
	return a, err
}

func (r *adminRepo) FindById(ctx context.Context, id int64) (*entity.AdminUser, error) {
	var a *entity.AdminUser
	err := g.Model("admin_user").Ctx(ctx).Where("id", id).Scan(&a)
	return a, err
}

func (r *adminRepo) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	_, err := g.Model("admin_user").Ctx(ctx).Where("id", id).Data(g.Map{
		"last_login_at": gtime.Now(),
		"last_ip":       ip,
	}).Update()
	return err
}

// ---------- 角色 ----------

func (r *adminRepo) ListRoles(ctx context.Context) ([]*entity.AdminRole, error) {
	var list []*entity.AdminRole
	err := g.Model("admin_role").Ctx(ctx).Order("id asc").Scan(&list)
	return list, err
}

func (r *adminRepo) FindRoleById(ctx context.Context, id int64) (*entity.AdminRole, error) {
	var role *entity.AdminRole
	err := g.Model("admin_role").Ctx(ctx).Where("id", id).Scan(&role)
	return role, err
}

func (r *adminRepo) FindRoleByCode(ctx context.Context, code string) (*entity.AdminRole, error) {
	var role *entity.AdminRole
	err := g.Model("admin_role").Ctx(ctx).Where("code", code).Scan(&role)
	return role, err
}

func (r *adminRepo) CreateRole(ctx context.Context, role *entity.AdminRole) (int64, error) {
	res, err := g.Model("admin_role").Ctx(ctx).Data(g.Map{
		"name":        role.Name,
		"code":        role.Code,
		"remark":      role.Remark,
		"status":      1,
		"permissions": role.Permissions,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *adminRepo) UpdateRole(ctx context.Context, id int64, name, remark string, status int, permissions string) error {
	_, err := g.Model("admin_role").Ctx(ctx).Where("id", id).Data(g.Map{
		"name":        name,
		"remark":      remark,
		"status":      status,
		"permissions": permissions,
		"updated_at":  gtime.Now(),
	}).Update()
	return err
}

func (r *adminRepo) DeleteRole(ctx context.Context, id int64) error {
	_, err := g.Model("admin_role").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *adminRepo) CountAdminsByRoleId(ctx context.Context, roleId int64) (int, error) {
	return g.Model("admin_user").Ctx(ctx).Where("role_id", roleId).Count()
}

// ---------- 管理员账号 ----------

func (r *adminRepo) ListAdmins(ctx context.Context, page, size int) ([]*entity.AdminUser, int, error) {
	m := g.Model("admin_user").Ctx(ctx)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.AdminUser
	err = m.Order("id asc").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *adminRepo) CreateAdmin(ctx context.Context, a *entity.AdminUser) (int64, error) {
	res, err := g.Model("admin_user").Ctx(ctx).Data(g.Map{
		"username": a.Username,
		"password": a.Password,
		"salt":     a.Salt,
		"nickname": a.Nickname,
		"role_id":  a.RoleId,
		"status":   1,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *adminRepo) UpdateAdmin(ctx context.Context, id int64, nickname string, roleId int64, status int, password, salt string) error {
	data := g.Map{
		"nickname":   nickname,
		"role_id":    roleId,
		"status":     status,
		"updated_at": gtime.Now(),
	}
	if password != "" {
		data["password"] = password
		data["salt"] = salt
	}
	_, err := g.Model("admin_user").Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

func (r *adminRepo) DeleteAdmin(ctx context.Context, id int64) error {
	_, err := g.Model("admin_user").Ctx(ctx).Where("id", id).Delete()
	return err
}
