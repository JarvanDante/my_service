// Package v1 后台管理员接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type AdminInfo struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoleId   int64  `json:"role_id"`
}

// 管理员登录(公开)
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"Backend/Auth" summary:"管理员登录"`
	Username string `json:"username" v:"required#账号必填"`
	Password string `json:"password" v:"required#密码必填"`
}
type LoginRes struct {
	Token string    `json:"token"`
	Admin AdminInfo `json:"admin"`
}

// 退出(需登录)
type LogoutReq struct {
	g.Meta `path:"/auth/logout" method:"post" tags:"Backend/Auth" summary:"管理员退出"`
}
type LogoutRes struct{}

// 当前管理员信息(需登录)
type InfoReq struct {
	g.Meta `path:"/auth/info" method:"get" tags:"Backend/Auth" summary:"当前管理员信息"`
}
type InfoRes struct {
	AdminInfo
}
