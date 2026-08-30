package v1

import "github.com/gogf/gf/v2/frame/g"

type ModuleItem struct {
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	Position      string   `json:"position"`
	Style         int      `json:"style"`
	Icon          int      `json:"icon"`
	CategoryIds   []int64  `json:"category_ids"`
	CategoryNames []string `json:"category_names"`
	TagIds        []int64  `json:"tag_ids"`
	TagNames      []string `json:"tag_names"`
	Size          int      `json:"size"`
	Rank          int      `json:"rank"`
	Status        int      `json:"status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type ModuleListReq struct {
	g.Meta     `path:"/comics-modules" method:"get" tags:"Backend/Comics" summary:"漫画模块列表"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	CategoryId int64  `json:"category_id"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}
type ModuleListRes struct {
	List  []ModuleItem `json:"list"`
	Total int          `json:"total"`
}

type ModuleCreateReq struct {
	g.Meta      `path:"/comics-modules" method:"post" tags:"Backend/Comics" summary:"新增漫画模块"`
	Name        string  `json:"name" v:"required#模块名必填"`
	Position    string  `json:"position"`
	Style       int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon        int     `json:"icon" v:"in:1,2,3#图标非法"`
	CategoryIds []int64 `json:"category_ids"`
	TagIds      []int64 `json:"tag_ids"`
	Size        int     `json:"size"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type ModuleCreateRes struct {
	Id int64 `json:"id"`
}

type ModuleUpdateReq struct {
	g.Meta      `path:"/comics-modules/{id}" method:"put" tags:"Backend/Comics" summary:"编辑漫画模块"`
	Id          int64   `json:"id" in:"path" v:"required|min:1#模块ID必填"`
	Name        string  `json:"name"`
	Position    string  `json:"position"`
	Style       int     `json:"style" v:"in:1,2,3,4,5,6,7#样式非法"`
	Icon        int     `json:"icon" v:"in:1,2,3#图标非法"`
	CategoryIds []int64 `json:"category_ids"`
	TagIds      []int64 `json:"tag_ids"`
	Size        int     `json:"size"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type ModuleUpdateRes struct{}

type ModuleDeleteReq struct {
	g.Meta `path:"/comics-modules/{id}" method:"delete" tags:"Backend/Comics" summary:"删除漫画模块"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#模块ID必填"`
}
type ModuleDeleteRes struct{}
