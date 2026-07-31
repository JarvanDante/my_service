// Package v1 后台管理员账号管理接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type AdminItem struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	RoleId      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	Status      int    `json:"status"`
	LastLoginAt string `json:"last_login_at"`
}

// 管理员列表
type AdminListReq struct {
	g.Meta `path:"/admins" method:"get" tags:"Backend/Admin" summary:"管理员列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type AdminListRes struct {
	List  []AdminItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// 创建管理员
type AdminCreateReq struct {
	g.Meta   `path:"/admins" method:"post" tags:"Backend/Admin" summary:"创建管理员"`
	Username string `json:"username" v:"required#账号必填"`
	Password string `json:"password" v:"required#密码必填"`
	Nickname string `json:"nickname"`
	RoleId   int64  `json:"role_id"  v:"required|min:1#角色必填"`
}
type AdminCreateRes struct {
	Id int64 `json:"id"`
}

// 更新管理员(password 为空表示不改)
type AdminUpdateReq struct {
	g.Meta   `path:"/admins/{id}" method:"put" tags:"Backend/Admin" summary:"更新管理员"`
	Id       int64  `json:"id"       v:"required|min:1#管理员ID必填|管理员ID必须大于0"`
	Nickname string `json:"nickname"`
	RoleId   int64  `json:"role_id"  v:"required|min:1#角色必填"`
	Status   int    `json:"status"`
	Password string `json:"password"`
}
type AdminUpdateRes struct{}

// 删除管理员
type AdminDeleteReq struct {
	g.Meta `path:"/admins/{id}" method:"delete" tags:"Backend/Admin" summary:"删除管理员"`
	Id     int64 `json:"id" v:"required|min:1#管理员ID必填|管理员ID必须大于0"`
}
type AdminDeleteRes struct{}
