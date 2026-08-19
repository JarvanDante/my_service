package v1

import "github.com/gogf/gf/v2/frame/g"

type FrontCategoryItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Kind int    `json:"kind"`
}

// CategoryListReq 启用中的漫画分类, 按权重倒序(公开)。
type CategoryListReq struct {
	g.Meta `path:"/comics/categories" method:"get" tags:"Front/Comics" summary:"漫画分类"`
}
type CategoryListRes struct {
	List []FrontCategoryItem `json:"list"`
}
