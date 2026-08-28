package v1

import "github.com/gogf/gf/v2/frame/g"

type ModuleItem struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Position  string   `json:"position"`
	Style     int      `json:"style"`
	Icon      int      `json:"icon"`
	TagIds    []int64  `json:"tag_ids"`
	TagNames  []string `json:"tag_names"`
	Size      int      `json:"size"`
	Rank      int      `json:"rank"`
	Status    int      `json:"status"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type VideoModuleListReq struct {
	g.Meta   `path:"/video-modules" method:"get" tags:"Backend/Video" summary:"视频模块列表"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type VideoModuleListRes struct {
	List  []ModuleItem `json:"list"`
	Total int          `json:"total"`
}

type VideoModuleCreateReq struct {
	g.Meta   `path:"/video-modules" method:"post" tags:"Backend/Video" summary:"新增视频模块"`
	Name     string  `json:"name" v:"required#模块名必填"`
	Position string  `json:"position"`
	Style    int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon     int     `json:"icon" v:"in:1,2,3#图标非法"`
	TagIds   []int64 `json:"tag_ids"`
	Size     int     `json:"size"`
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type VideoModuleCreateRes struct {
	Id int64 `json:"id"`
}

type VideoModuleUpdateReq struct {
	g.Meta   `path:"/video-modules/{id}" method:"put" tags:"Backend/Video" summary:"编辑视频模块"`
	Id       int64   `json:"id" in:"path" v:"required|min:1#模块ID必填"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Style    int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon     int     `json:"icon" v:"in:1,2,3#图标非法"`
	TagIds   []int64 `json:"tag_ids"`
	Size     int     `json:"size"`
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type VideoModuleUpdateRes struct{}

type VideoModuleDeleteReq struct {
	g.Meta `path:"/video-modules/{id}" method:"delete" tags:"Backend/Video" summary:"删除视频模块"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#模块ID必填"`
}
type VideoModuleDeleteRes struct{}

type CartoonModuleListReq struct {
	g.Meta   `path:"/cartoon-modules" method:"get" tags:"Backend/Cartoon" summary:"动漫模块列表"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type CartoonModuleListRes struct {
	List  []ModuleItem `json:"list"`
	Total int          `json:"total"`
}

type CartoonModuleCreateReq struct {
	g.Meta   `path:"/cartoon-modules" method:"post" tags:"Backend/Cartoon" summary:"新增动漫模块"`
	Name     string  `json:"name" v:"required#模块名必填"`
	Position string  `json:"position"`
	Style    int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon     int     `json:"icon" v:"in:1,2,3#图标非法"`
	TagIds   []int64 `json:"tag_ids"`
	Size     int     `json:"size"`
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type CartoonModuleCreateRes struct {
	Id int64 `json:"id"`
}

type CartoonModuleUpdateReq struct {
	g.Meta   `path:"/cartoon-modules/{id}" method:"put" tags:"Backend/Cartoon" summary:"编辑动漫模块"`
	Id       int64   `json:"id" in:"path" v:"required|min:1#模块ID必填"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Style    int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon     int     `json:"icon" v:"in:1,2,3#图标非法"`
	TagIds   []int64 `json:"tag_ids"`
	Size     int     `json:"size"`
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type CartoonModuleUpdateRes struct{}

type CartoonModuleDeleteReq struct {
	g.Meta `path:"/cartoon-modules/{id}" method:"delete" tags:"Backend/Cartoon" summary:"删除动漫模块"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#模块ID必填"`
}
type CartoonModuleDeleteRes struct{}
