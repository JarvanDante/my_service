// Package v1 前台收藏/点赞契约(移植自 tianbi collect)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// OperateReq 添加/取消 收藏/点赞/踩(id 与 ids 二选一, ids 为同类型批量)。
type OperateReq struct {
	g.Meta    `path:"/collect/operate" method:"post" tags:"Front/Collect" summary:"添加/取消收藏点赞"`
	Id        int64   `json:"id"`
	Ids       []int64 `json:"ids"`
	MediaType int     `json:"media_type" v:"required|in:1,2#资源类型必填|资源类型非法"` // 1视频 2帖子
	Flag      bool    `json:"flag"`                                         // true=添加 false=取消
	Type      int     `json:"type" v:"required|in:1,2,3#操作类型必填|操作类型非法"`     // 1收藏 2点赞 3踩
}
type OperateRes struct{}

// DeleteReq 批量取消(移植自 tianbi collect/delete)。
type DeleteReq struct {
	g.Meta    `path:"/collect/delete" method:"post" tags:"Front/Collect" summary:"批量取消收藏"`
	Ids       []int64 `json:"ids" v:"required#ids必填"`
	MediaType int     `json:"media_type" v:"required|in:1,2#资源类型必填|资源类型非法"`
	Type      int     `json:"type" v:"required|in:1,2,3#操作类型必填|操作类型非法"`
}
type DeleteRes struct{}

type CollectItem struct {
	ContentId int64  `json:"content_id"`
	MediaType int    `json:"media_type"`
	CreatedAt string `json:"created_at"`
}

// ListReq 我的收藏/点赞列表(返回 content_id, 客户端按需拉详情)。
type ListReq struct {
	g.Meta    `path:"/collect/list" method:"get" tags:"Front/Collect" summary:"我的收藏/点赞列表"`
	Type      int `json:"type" v:"required|in:1,2,3#操作类型必填|操作类型非法"` // 1收藏 2点赞 3踩
	MediaType int `json:"media_type"`                               // 0=全部
	Page      int `json:"page"`
	Size      int `json:"size"`
}
type ListRes struct {
	List  []CollectItem `json:"list"`
	Total int           `json:"total"`
}
