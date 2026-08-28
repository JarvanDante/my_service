package v1

import "github.com/gogf/gf/v2/frame/g"

type FrontModuleItem struct {
	Id    int64    `json:"id"`
	Name  string   `json:"name"`
	Style int      `json:"style"`
	Icon  int      `json:"icon"`
	Size  int      `json:"size"`
	Tags  []string `json:"tags"`
	Items []Item   `json:"items"`
}

type VideoModuleListReq struct {
	g.Meta   `path:"/video/modules" method:"get" tags:"Front/Video" summary:"视频首页模块"`
	Position string `json:"position"`
}
type VideoModuleListRes struct {
	List []FrontModuleItem `json:"list"`
}

type CartoonModuleListReq struct {
	g.Meta   `path:"/cartoon/modules" method:"get" tags:"Front/Cartoon" summary:"动漫首页模块"`
	Position string `json:"position"`
}
type CartoonModuleListRes struct {
	List []FrontModuleItem `json:"list"`
}
