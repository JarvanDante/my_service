package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id             int64  `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Remark         string `json:"remark"`
	Status         int    `json:"status"`
	StatusText     string `json:"status_text"`
	UserCount      int    `json:"user_count"`
	H5Link         string `json:"h5_link"`
	ClipboardValue string `json:"clipboard_value"`
	CreatedAt      string `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/channels" method:"get" tags:"Backend/Channel" summary:"渠道列表"`
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
	g.Meta `path:"/channels" method:"post" tags:"Backend/Channel" summary:"新增渠道"`
	Code   string `json:"code" v:"required#渠道码必填"`
	Name   string `json:"name" v:"required#渠道名必填"`
	Remark string `json:"remark"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/channels/{id}" method:"put" tags:"Backend/Channel" summary:"编辑渠道"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name   string `json:"name" v:"required#渠道名必填"`
	Remark string `json:"remark"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/channels/{id}" method:"delete" tags:"Backend/Channel" summary:"删除渠道"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
