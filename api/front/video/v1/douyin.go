package v1

import "github.com/gogf/gf/v2/frame/g"

// DouyinCategoryListReq 启用中的抖音分类, 按权重倒序(公开)。
type DouyinCategoryListReq struct {
	g.Meta `path:"/douyin/categories" method:"get" tags:"Front/Douyin" summary:"抖音分类"`
}
type DouyinCategoryListRes struct {
	List []FrontCategoryItem `json:"list"`
}

// DouyinListReq 抖音列表(公开, video.kind=3 且已上架)。
// sort: 0综合(人工 sort 权重) 1最新 2时长。
type DouyinListReq struct {
	g.Meta   `path:"/douyin/list" method:"get" tags:"Front/Douyin" summary:"抖音列表"`
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
	Sort     int    `json:"sort" v:"in:0,1,2#排序方式非法"`
	Follow   int    `json:"follow"` // 1=只看已关注的UP主
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type DouyinListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}
