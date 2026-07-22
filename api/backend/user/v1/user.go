// Package v1 后台用户接口契约(面向运营/管理后台, 含写操作)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type DisableReq struct {
	g.Meta `path:"/users/{id}/disable" method:"post" tags:"Backend/User" summary:"禁用用户(后台)"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填"`
}
type DisableRes struct{}

type DetailReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"Backend/User" summary:"用户详情(后台)"`
	Id     int64 `json:"id" v:"required|min:1"`
}
type DetailRes struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
