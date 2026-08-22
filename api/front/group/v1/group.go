package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Intro    string `json:"intro"`
	Avatar   string `json:"avatar"`
	Link     string `json:"link"`
	Platform string `json:"platform"`
}

type ListReq struct {
	g.Meta `path:"/group/list" method:"get" tags:"Front/Group" summary:"官方社群列表"`
}
type ListRes struct {
	List []Item `json:"list"`
}
