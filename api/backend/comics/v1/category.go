package v1

import "github.com/gogf/gf/v2/frame/g"

type CategoryItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Kind      int    `json:"kind"`
	Rank      int    `json:"rank"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type CategoryListReq struct {
	g.Meta `path:"/comics-categories" method:"get" tags:"Backend/Comics" summary:"漫画分类列表"`
	Kind   string `json:"kind"`   // 空=全部  0/1/2/3
	Status string `json:"status"` // 空=全部  0禁用  1启用
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type CategoryListRes struct {
	List  []CategoryItem `json:"list"`
	Total int            `json:"total"`
}

type CategoryCreateReq struct {
	g.Meta `path:"/comics-categories" method:"post" tags:"Backend/Comics" summary:"新增漫画分类"`
	Name   string `json:"name" v:"required#分类名必填"`
	Kind   int    `json:"kind" v:"in:0,1,2,3#类型非法"`
	Rank   int    `json:"rank"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CategoryCreateRes struct {
	Id int64 `json:"id"`
}

type CategoryUpdateReq struct {
	g.Meta `path:"/comics-categories/{id}" method:"put" tags:"Backend/Comics" summary:"编辑漫画分类"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#分类ID必填"`
	Name   string `json:"name"`
	Kind   int    `json:"kind" v:"in:0,1,2,3#类型非法"`
	Rank   int    `json:"rank"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CategoryUpdateRes struct{}

type CategoryDeleteReq struct {
	g.Meta `path:"/comics-categories/{id}" method:"delete" tags:"Backend/Comics" summary:"删除漫画分类"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#分类ID必填"`
}
type CategoryDeleteRes struct{}
