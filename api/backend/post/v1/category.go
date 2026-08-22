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
	g.Meta `path:"/post-categories" method:"get" tags:"Backend/Post" summary:"帖子分类列表"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type CategoryListRes struct {
	List  []CategoryItem `json:"list"`
	Total int            `json:"total"`
}

type CategoryCreateReq struct {
	g.Meta `path:"/post-categories" method:"post" tags:"Backend/Post" summary:"新增帖子分类"`
	Name   string `json:"name" v:"required#分类名必填"`
	Kind   int    `json:"kind" v:"in:0,1,2#类型非法"`
	Rank   int    `json:"rank"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CategoryCreateRes struct {
	Id int64 `json:"id"`
}

type CategoryUpdateReq struct {
	g.Meta `path:"/post-categories/{id}" method:"put" tags:"Backend/Post" summary:"编辑帖子分类"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#分类ID必填"`
	Name   string `json:"name"`
	Kind   int    `json:"kind" v:"in:0,1,2#类型非法"`
	Rank   int    `json:"rank"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CategoryUpdateRes struct{}

type CategoryDeleteReq struct {
	g.Meta `path:"/post-categories/{id}" method:"delete" tags:"Backend/Post" summary:"删除帖子分类"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#分类ID必填"`
}
type CategoryDeleteRes struct{}
