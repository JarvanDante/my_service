// Package service 后台管理员对外接口。
package service

import "context"

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

// RoleDTO 角色。
type RoleDTO struct {
	Id     int64
	Name   string
	Code   string
	Remark string
	Status int
}

// PermDTO 一条权限(路径 + 方法)。
type PermDTO struct {
	Path   string
	Method string
}

type IAdmin interface {
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	Logout(ctx context.Context, adminId int64) error
	Info(ctx context.Context, adminId int64) (*AdminInfoDTO, error)

	// 角色 / 权限管理
	ListRoles(ctx context.Context) ([]*RoleDTO, error)
	ListPerms(ctx context.Context, roleCode string) ([]*PermDTO, error)
	AddPerm(ctx context.Context, roleCode, path, method string) error
	RemovePerm(ctx context.Context, roleCode, path, method string) error
}
