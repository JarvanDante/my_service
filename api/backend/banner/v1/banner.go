package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64  `json:"id"`
	Position  string `json:"position"`
	Title     string `json:"title"`
	CoverUrl  string `json:"cover_url"`
	Link      string `json:"link"`
	Rank      int    `json:"rank"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListReq struct {
	g.Meta   `path:"/banner" method:"get" tags:"Backend/Banner" summary:"轮播列表"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta   `path:"/banner" method:"post" tags:"Backend/Banner" summary:"新增轮播"`
	Position string `json:"position" v:"required#位置必填"`
	Title    string `json:"title"`
	CoverUrl string `json:"cover_url" v:"required#请上传轮播图"`
	Link     string `json:"link"`
	Rank     int    `json:"rank"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/banner/{id}" method:"put" tags:"Backend/Banner" summary:"更新轮播"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Position string `json:"position" v:"required#位置必填"`
	Title    string `json:"title"`
	CoverUrl string `json:"cover_url"`
	Link     string `json:"link"`
	Rank     int    `json:"rank"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/banner/{id}" method:"delete" tags:"Backend/Banner" summary:"删除轮播"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
