// Package domain 用户领域层。纯 Go, 禁止 import 任何 github.com/gogf/gf。
package domain

import "errors"

type Status int

const (
	StatusActive   Status = 1
	StatusDisabled Status = 2
)

// User 领域实体, 业务规则内聚在方法上。
type User struct {
	ID     int64
	Name   string
	Email  string
	Status Status
}

func (u *User) IsActive() bool { return u.Status == StatusActive }

// Disable 领域规则: 已禁用不可重复禁用。
func (u *User) Disable() error {
	if u.Status == StatusDisabled {
		return errors.New("user already disabled")
	}
	u.Status = StatusDisabled
	return nil
}
