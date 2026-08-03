// Code generated and maintained manually (RBAC 菜单+接口权限树).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type AdminPermission struct {
	Id         int64       `json:"id"         orm:"id"`
	ParentId   int64       `json:"parentId"   orm:"parent_id"`
	Name       string      `json:"name"       orm:"name"`
	RouteUrl   string      `json:"routeUrl"   orm:"route_url"`
	Component  string      `json:"component"  orm:"component"`
	Method     string      `json:"method"     orm:"method"`
	Icon       string      `json:"icon"       orm:"icon"`
	IsMenu     int         `json:"isMenu"     orm:"is_menu"`
	HideInMenu int         `json:"hideInMenu" orm:"hide_in_menu"`
	AffixTab   int         `json:"affixTab"   orm:"affix_tab"`
	ActivePath string      `json:"activePath" orm:"active_path"`
	Sort       int         `json:"sort"       orm:"sort"`
	Status     int         `json:"status"     orm:"status"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}
