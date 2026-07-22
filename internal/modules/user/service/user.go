// Package service 用户模块对外能力接口。其他模块只能依赖此接口。
package service

import "context"

// UserDTO 对外数据传输对象, 不外泄领域实体。
type UserDTO struct {
	ID     int64
	Name   string
	Active bool
}

type IUser interface {
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
	DisableUser(ctx context.Context, id int64) error
}
