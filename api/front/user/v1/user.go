// Package v1 前台用户接口契约(面向 C 端, 只读为主)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type GetReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"Front/User" summary:"获取用户(前台)"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填"`
}
type GetRes struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
