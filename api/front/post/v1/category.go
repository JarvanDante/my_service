package v1

import "github.com/gogf/gf/v2/frame/g"

type FrontCategoryItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Kind int    `json:"kind"` // 0普通 1最新 2推荐 3榜单
}

type CategoryListReq struct {
	g.Meta `path:"/post/categories" method:"get" tags:"Front/Post" summary:"帖子分类"`
}
type CategoryListRes struct {
	List []FrontCategoryItem `json:"list"`
}
