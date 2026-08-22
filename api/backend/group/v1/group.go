package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	Avatar    string `json:"avatar"`
	Link      string `json:"link"`
	Platform  string `json:"platform"`
	Rank      int    `json:"rank"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/group" method:"get" tags:"Backend/Group" summary:"官方社群列表"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta   `path:"/group" method:"post" tags:"Backend/Group" summary:"新增官方社群"`
	Name     string `json:"name" v:"required#社群名必填"`
	Intro    string `json:"intro"`
	Avatar   string `json:"avatar"`
	Link     string `json:"link" v:"required#跳转链接必填"`
	Platform string `json:"platform"`
	Rank     int    `json:"rank"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/group/{id}" method:"put" tags:"Backend/Group" summary:"更新官方社群"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name     string `json:"name"`
	Intro    string `json:"intro"`
	Avatar   string `json:"avatar"`
	Link     string `json:"link"`
	Platform string `json:"platform"`
	Rank     int    `json:"rank"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/group/{id}" method:"delete" tags:"Backend/Group" summary:"删除官方社群"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
