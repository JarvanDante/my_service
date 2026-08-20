package v1

import "github.com/gogf/gf/v2/frame/g"

type FrontCategoryItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Kind int    `json:"kind"` // 0普通 1最新 2推荐 3榜单
}

// CartoonCategoryListReq 启用中的动漫分类, 按权重倒序(公开)。
type CartoonCategoryListReq struct {
	g.Meta `path:"/cartoon/categories" method:"get" tags:"Front/Cartoon" summary:"动漫分类"`
}
type CartoonCategoryListRes struct {
	List []FrontCategoryItem `json:"list"`
}

// CartoonListReq 动漫列表(公开, video.kind=2 且已上架)。
// sort: 0综合(人工 sort 权重) 1最新 2时长。
type CartoonListReq struct {
	g.Meta   `path:"/cartoon/list" method:"get" tags:"Front/Cartoon" summary:"动漫列表"`
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Sort     int    `json:"sort" v:"in:0,1,2#排序方式非法"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type CartoonListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}
