// Package v1 子后台 RBAC 菜单+接口权限树接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type PermNode struct {
	Id         int64      `json:"id"`
	ParentId   int64      `json:"parentId"`
	Name       string     `json:"name"`
	RouteUrl   string     `json:"routeUrl"`
	Component  string     `json:"component"`
	Method     string     `json:"method"`
	Icon       string     `json:"icon"`
	IsMenu     int        `json:"isMenu"`
	HideInMenu int        `json:"hideInMenu"`
	AffixTab   int        `json:"affixTab"`
	ActivePath string     `json:"activePath"`
	Sort       int        `json:"sort"`
	Status     int        `json:"status"`
	Children   []PermNode `json:"children"`
}

type PermTreeReq struct {
	g.Meta `path:"/permissions" method:"get" tags:"Backend/Permission" summary:"权限树(菜单+接口)"`
}
type PermTreeRes struct {
	List []PermNode `json:"list"`
}

type PermCreateReq struct {
	g.Meta     `path:"/permissions" method:"post" tags:"Backend/Permission" summary:"新增权限节点"`
	ParentId   int64  `json:"parentId"`
	Name       string `json:"name" v:"required#名称必填"`
	RouteUrl   string `json:"routeUrl"`
	Component  string `json:"component"`
	Method     string `json:"method"`
	Icon       string `json:"icon"`
	IsMenu     int    `json:"isMenu"`
	HideInMenu int    `json:"hideInMenu"`
	AffixTab   int    `json:"affixTab"`
	ActivePath string `json:"activePath"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
}
type PermCreateRes struct {
	Id int64 `json:"id"`
}

type PermUpdateReq struct {
	g.Meta     `path:"/permissions/{id}" method:"put" tags:"Backend/Permission" summary:"更新权限节点"`
	Id         int64  `json:"id" v:"required|min:1#ID必填|ID必须大于0"`
	ParentId   int64  `json:"parentId"`
	Name       string `json:"name" v:"required#名称必填"`
	RouteUrl   string `json:"routeUrl"`
	Component  string `json:"component"`
	Method     string `json:"method"`
	Icon       string `json:"icon"`
	IsMenu     int    `json:"isMenu"`
	HideInMenu int    `json:"hideInMenu"`
	AffixTab   int    `json:"affixTab"`
	ActivePath string `json:"activePath"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
}
type PermUpdateRes struct{}

type PermDeleteReq struct {
	g.Meta `path:"/permissions/{id}" method:"delete" tags:"Backend/Permission" summary:"删除权限节点"`
	Id     int64 `json:"id" v:"required|min:1#ID必填|ID必须大于0"`
}
type PermDeleteRes struct{}
