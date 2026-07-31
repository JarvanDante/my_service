// Package v1 后台角色 / 权限接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type RoleItem struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

// 角色列表
type ListRolesReq struct {
	g.Meta `path:"/roles" method:"get" tags:"Backend/Role" summary:"角色列表"`
}
type ListRolesRes struct {
	List []RoleItem `json:"list"`
}

type PermItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// 角色权限列表
type ListPermsReq struct {
	g.Meta `path:"/roles/{code}/perms" method:"get" tags:"Backend/Role" summary:"角色权限列表"`
	Code   string `json:"code" v:"required#角色码必填"`
}
type ListPermsRes struct {
	List []PermItem `json:"list"`
}

// 新增角色权限
type AddPermReq struct {
	g.Meta `path:"/roles/{code}/perms" method:"post" tags:"Backend/Role" summary:"新增角色权限"`
	Code   string `json:"code"   v:"required#角色码必填"`
	Path   string `json:"path"   v:"required#路径必填"`
	Method string `json:"method" v:"required#方法必填"`
}
type AddPermRes struct{}

// 删除角色权限
type DelPermReq struct {
	g.Meta `path:"/roles/{code}/perms" method:"delete" tags:"Backend/Role" summary:"删除角色权限"`
	Code   string `json:"code"   v:"required#角色码必填"`
	Path   string `json:"path"   v:"required#路径必填"`
	Method string `json:"method" v:"required#方法必填"`
}
type DelPermRes struct{}

// ---------- 角色管理(CRUD) ----------

// 创建角色
type RoleCreateReq struct {
	g.Meta `path:"/roles" method:"post" tags:"Backend/Role" summary:"创建角色"`
	Name   string `json:"name"   v:"required#角色名必填"`
	Code   string `json:"code"   v:"required#角色码必填"`
	Remark string `json:"remark"`
}
type RoleCreateRes struct {
	Id int64 `json:"id"`
}

// 更新角色
type RoleUpdateReq struct {
	g.Meta `path:"/roles/{id}" method:"put" tags:"Backend/Role" summary:"更新角色"`
	Id     int64  `json:"id"     v:"required|min:1#角色ID必填|角色ID必须大于0"`
	Name   string `json:"name"   v:"required#角色名必填"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}
type RoleUpdateRes struct{}

// 删除角色
type RoleDeleteReq struct {
	g.Meta `path:"/roles/{id}" method:"delete" tags:"Backend/Role" summary:"删除角色"`
	Id     int64 `json:"id" v:"required|min:1#角色ID必填|角色ID必须大于0"`
}
type RoleDeleteRes struct{}
