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

// ModuleListReq 启用中的漫画分类模块, 按权重倒序(公开)。position 空则用权重最高的分类。
type ModuleListReq struct {
	g.Meta   `path:"/comics/modules" method:"get" tags:"Front/Comics" summary:"漫画分类模块"`
	Position string `json:"position"`
}
type ModuleListRes struct {
	List []FrontModuleItem `json:"list"`
}

type ModuleRefreshReq struct {
	g.Meta  `path:"/comics/modules/refresh" method:"get" tags:"Front/Comics" summary:"漫画模块换一换"`
	Id      int64  `json:"id" v:"required|min:1"`
	Exclude string `json:"exclude"`
}
type ModuleRefreshRes struct {
	FrontModuleItem
}
