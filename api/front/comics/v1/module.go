package v1

import "github.com/gogf/gf/v2/frame/g"

type FrontModuleItem struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	Style      int      `json:"style"`
	Icon       int      `json:"icon"`
	Size       int      `json:"size"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Items      []Item   `json:"items"`
}

// ModuleListReq 启用中的漫画首页模块, 按权重倒序(公开)。
type ModuleListReq struct {
	g.Meta   `path:"/comics/modules" method:"get" tags:"Front/Comics" summary:"漫画首页模块"`
	Position string `json:"position"`
}
type ModuleListRes struct {
	List []FrontModuleItem `json:"list"`
}
