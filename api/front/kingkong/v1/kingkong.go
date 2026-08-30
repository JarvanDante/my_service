package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IconUrl  string `json:"icon_url"`
	OpenMode string `json:"open_mode"`
	Link     string `json:"link"`
	AppLink  string `json:"app_link"`
	Position string `json:"position"`
	Sort     int    `json:"sort"`
}

type ListReq struct {
	g.Meta   `path:"/kingkong/list" method:"get" tags:"Front/Kingkong" summary:"金刚区列表"`
	Position string `json:"position" v:"required#展示位置必填"`
}
type ListRes struct {
	List []Item `json:"list"`
}
