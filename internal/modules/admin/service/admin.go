// Package service 后台管理员对外接口。
package service

import "context"

// ---------- 认证 ----------

type LoginInput struct {
	Username string
	Password string
	Ip       string
}

type AdminInfoDTO struct {
	Id       int64
	Username string
	Nickname string
	RoleId   int64
}

type LoginDTO struct {
	Token string
	Admin *AdminInfoDTO
}

// ---------- 角色 / 权限 ----------

type RoleDTO struct {
	Id     int64
	Name   string
	Code   string
	Remark string
	Status int
}

type PermDTO struct {
	Path   string
	Method string
}

type RoleCreateInput struct {
	Name   string
	Code   string
	Remark string
}

type RoleUpdateInput struct {
	Id     int64
	Name   string
	Remark string
	Status int
}

// ---------- 管理员账号 ----------

type AdminItemDTO struct {
	Id          int64
	Username    string
	Nickname    string
	RoleId      int64
	RoleName    string
	Status      int
	LastLoginAt string
}

type AdminListDTO struct {
	List  []*AdminItemDTO
	Total int
	Page  int
	Size  int
}

type AdminCreateInput struct {
	Username string
	Password string
	Nickname string
	RoleId   int64
}

type AdminUpdateInput struct {
	Id       int64
	Nickname string
	RoleId   int64
	Status   int
	Password string // 空表示不改
}

type IAdmin interface {
	// 认证
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	Logout(ctx context.Context, adminId int64) error
	Info(ctx context.Context, adminId int64) (*AdminInfoDTO, error)

	// 角色权限(Casbin)
	ListRoles(ctx context.Context) ([]*RoleDTO, error)
	ListPerms(ctx context.Context, roleCode string) ([]*PermDTO, error)
	AddPerm(ctx context.Context, roleCode, path, method string) error
	RemovePerm(ctx context.Context, roleCode, path, method string) error

	// 角色管理
	CreateRole(ctx context.Context, in RoleCreateInput) (int64, error)
	UpdateRole(ctx context.Context, in RoleUpdateInput) error
	DeleteRole(ctx context.Context, id int64) error

	// 管理员账号管理
	ListAdmins(ctx context.Context, page, size int) (*AdminListDTO, error)
	CreateAdmin(ctx context.Context, in AdminCreateInput) (int64, error)
	UpdateAdmin(ctx context.Context, in AdminUpdateInput) error
	DeleteAdmin(ctx context.Context, id, operatorId int64) error
}
