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

type IAdmin interface {
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	Logout(ctx context.Context, adminId int64) error
	Info(ctx context.Context, adminId int64) (*AdminInfoDTO, error)
}
