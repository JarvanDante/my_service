// Package v1 后台热搜词契约(CRUD + 排行缓存刷新)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id          int64  `json:"id"`
	Keyword     string `json:"keyword"`
	Heat        int    `json:"heat"`
	SearchCount int64  `json:"search_count"`
	Status      int    `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

type ListReq struct {
	g.Meta  `path:"/hotsearch" method:"get" tags:"Backend/Rank" summary:"热搜词列表"`
	Status  string `json:"status"` // 空=全部
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta  `path:"/hotsearch" method:"post" tags:"Backend/Rank" summary:"新增热搜词"`
	Keyword string `json:"keyword" v:"required#关键词必填"`
	Heat    int    `json:"heat"`
	Status  int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta  `path:"/hotsearch/{id}" method:"put" tags:"Backend/Rank" summary:"更新热搜词"`
	Id      int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Keyword string `json:"keyword"`
	Heat    int    `json:"heat"`
	Status  int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/hotsearch/{id}" method:"delete" tags:"Backend/Rank" summary:"删除热搜词"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// RefreshRankReq 清除排行缓存(下次请求重算)。
type RefreshRankReq struct {
	g.Meta `path:"/rank/refresh" method:"post" tags:"Backend/Rank" summary:"刷新排行缓存"`
}
type RefreshRankRes struct{}
