package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	IconUrl      string `json:"icon_url"`
	OpenMode     string `json:"open_mode"`
	OpenModeName string `json:"open_mode_name"`
	Link         string `json:"link"`
	AppLink      string `json:"app_link"`
	LinkLabel    string `json:"link_label"`
	Position     string `json:"position"`
	PositionName string `json:"position_name"`
	Sort         int    `json:"sort"`
	Status       int    `json:"status"`
	StatusText   string `json:"status_text"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ListReq struct {
	g.Meta   `path:"/kingkong-items" method:"get" tags:"Backend/Kingkong" summary:"金刚区列表"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta   `path:"/kingkong-items" method:"post" tags:"Backend/Kingkong" summary:"新增金刚区"`
	Name     string `json:"name" v:"required#名称必填"`
	IconUrl  string `json:"icon_url" v:"required#请上传图标"`
	OpenMode string `json:"open_mode" v:"required|in:block,list,douyin#打开方式非法"`
	Link     string `json:"link"`
	AppLink  string `json:"app_link"`
	Position string `json:"position" v:"required|in:comics,cartoon,movie,novel,short#请选择展示位置"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/kingkong-items/{id}" method:"put" tags:"Backend/Kingkong" summary:"编辑金刚区"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name     string `json:"name" v:"required#名称必填"`
	IconUrl  string `json:"icon_url" v:"required#请上传图标"`
	OpenMode string `json:"open_mode" v:"required|in:block,list,douyin#打开方式非法"`
	Link     string `json:"link"`
	AppLink  string `json:"app_link"`
	Position string `json:"position" v:"required|in:comics,cartoon,movie,novel,short#请选择展示位置"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/kingkong-items/{id}" method:"delete" tags:"Backend/Kingkong" summary:"删除金刚区"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
