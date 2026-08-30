package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id       int64  `json:"id"`
	Position string `json:"position"`
	Title    string `json:"title"`
	CoverUrl string `json:"cover_url"`
	Link     string `json:"link"`
}

type ListReq struct {
	g.Meta   `path:"/banner/list" method:"get" tags:"Front/Banner" summary:"首页轮播"`
	Position string `json:"position" v:"required#位置必填"`
}
type ListRes struct {
	List []Item `json:"list"`
}
