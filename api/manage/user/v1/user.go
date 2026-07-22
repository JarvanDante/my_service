// Package v1 总后台用户接口契约(面向平台超管)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type GetReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"Manage/User" summary:"获取用户(总后台)"`
	Id     int64 `json:"id" v:"required|min:1"`
}
type GetRes struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
