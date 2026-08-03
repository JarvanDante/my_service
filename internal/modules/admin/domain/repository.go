// Package domain 后台管理员领域层。
package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	// 认证
	FindByUsername(ctx context.Context, username string) (*entity.AdminUser, error)
	FindById(ctx context.Context, id int64) (*entity.AdminUser, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error

	// 角色
	ListRoles(ctx context.Context) ([]*entity.AdminRole, error)
	FindRoleById(ctx context.Context, id int64) (*entity.AdminRole, error)
	FindRoleByCode(ctx context.Context, code string) (*entity.AdminRole, error)
	CreateRole(ctx context.Context, r *entity.AdminRole) (int64, error)
	UpdateRole(ctx context.Context, id int64, name, remark string, status int, permissions string) error
	DeleteRole(ctx context.Context, id int64) error
	CountAdminsByRoleId(ctx context.Context, roleId int64) (int, error)

	// 权限树(RBAC 菜单+接口权限)
	ListPermissions(ctx context.Context) ([]*entity.AdminPermission, error)
	FindPermissionsByIds(ctx context.Context, ids []int64) ([]*entity.AdminPermission, error)
	FindPermissionById(ctx context.Context, id int64) (*entity.AdminPermission, error)
	CreatePermission(ctx context.Context, p *entity.AdminPermission) (int64, error)
	UpdatePermission(ctx context.Context, p *entity.AdminPermission) error
	DeletePermission(ctx context.Context, id int64) error
	CountPermissionChildren(ctx context.Context, parentId int64) (int, error)

	// 管理员账号
	ListAdmins(ctx context.Context, page, size int) ([]*entity.AdminUser, int, error)
	CreateAdmin(ctx context.Context, a *entity.AdminUser) (int64, error)
	UpdateAdmin(ctx context.Context, id int64, nickname string, roleId int64, status int, password, salt string) error
	DeleteAdmin(ctx context.Context, id int64) error
}
