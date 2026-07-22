// Package service 用户模块对外能力接口。其他模块/门面只依赖此接口。
package service

import "context"

// LoginInput 登录入参(设备登录)。
type LoginInput struct {
	DeviceId      string
	DeviceType    string
	DeviceVersion string
	Ip            string
}

// UserInfoDTO 前台用户详情。
type UserInfoDTO struct {
	Id        int64
	Username  string
	Nickname  string
	Phone     string
	Img       string
	Signature string
	Sex       int
	Level     int
	Balance   float64
	Credit    float64
	GroupName string
	Fans      int
	Follow    int
}

// LoginDTO 登录返回。
type LoginDTO struct {
	Token string
	User  *UserInfoDTO
}

// UserDTO 后台/总后台用的精简信息。
type UserDTO struct {
	ID     int64
	Name   string
	Active bool
}

type IUser interface {
	// Login 设备登录/游客自动注册。
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	// Info 当前登录用户详情。
	Info(ctx context.Context, userId int64) (*UserInfoDTO, error)
	// GetUser 精简信息(后台/总后台)。
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
	// DisableUser 禁用(后台)。
	DisableUser(ctx context.Context, id int64) error
}
