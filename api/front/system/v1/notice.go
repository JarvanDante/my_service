package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListReq struct {
	g.Meta `path:"/notice/list" method:"get" tags:"Front/System" summary:"上架中的系统公告"`
}
type ListRes struct {
	List []Item `json:"list"`
}
