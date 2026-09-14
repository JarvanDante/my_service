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

type ModuleListReq struct {
	g.Meta   `path:"/novel/modules" method:"get" tags:"Front/Novel" summary:"小说分类模块"`
	Position string `json:"position"`
}
type ModuleListRes struct {
	List []FrontModuleItem `json:"list"`
}

type ModuleRefreshReq struct {
	g.Meta  `path:"/novel/modules/refresh" method:"get" tags:"Front/Novel" summary:"小说模块换一换"`
	Id      int64  `json:"id" v:"required|min:1"`
	Exclude string `json:"exclude"`
}
type ModuleRefreshRes struct {
	FrontModuleItem
}
