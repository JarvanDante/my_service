// Package v1 子后台菜单下发接口契约(动态菜单 method B)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// MenusReq 当前管理员菜单树(需登录)。
type MenusReq struct {
	g.Meta `path:"/auth/menus" method:"get" tags:"Backend/Auth" summary:"当前管理员菜单树"`
}

// MenuMeta 对应 vben 路由 meta。
type MenuMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon,omitempty"`
	Order      int    `json:"order,omitempty"`
	AffixTab   bool   `json:"affixTab,omitempty"`
	HideInMenu bool   `json:"hideInMenu,omitempty"`
	ActivePath string `json:"activePath,omitempty"`
}

// MenuNode 对应 vben RouteRecordStringComponent。
type MenuNode struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component,omitempty"`
	Redirect  string     `json:"redirect,omitempty"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuNode `json:"children,omitempty"`
}

// MenusRes 菜单树列表。
type MenusRes struct {
	List []MenuNode `json:"list"`
}
