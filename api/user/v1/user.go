// Package v1 用户模块对外接口契约(req/res)。仅数据结构 + 校验 tag。
package v1

import "github.com/gogf/gf/v2/frame/g"

type GetReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"User" summary:"获取用户"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填"`
}
type GetRes struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type DisableReq struct {
	g.Meta `path:"/users/{id}/disable" method:"post" tags:"User" summary:"禁用用户"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填"`
}
type DisableRes struct{}
