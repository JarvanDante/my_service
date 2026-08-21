// Package v1 前台评论契约(移植自 tianbi comment)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	ParentId   int64  `json:"parent_id"`
	RootId     int64  `json:"root_id"`
	Content    string `json:"content"`
	LikeCount  int    `json:"like_count"`
	ReplyCount int    `json:"reply_count"`
	CreatedAt  string `json:"created_at"`
	Replies    []Item `json:"replies,omitempty"` // 顶层评论携带回复
}

// AddReq 发表评论/回复(需登录, 过敏感词)。
type AddReq struct {
	g.Meta    `path:"/comment/add" method:"post" tags:"Front/Comment" summary:"发表评论"`
	MediaType int    `json:"media_type" v:"required|in:1,2,4,7#资源类型必填|资源类型非法"`
	ContentId int64  `json:"content_id" v:"required|min:1#内容ID必填"`
	ParentId  int64  `json:"parent_id"` // 0=顶层, >0=回复
	Content   string `json:"content" v:"required|max-length:1000#内容必填|内容过长"`
}
type AddRes struct {
	Id     int64 `json:"id"`
	Status int   `json:"status"` // 0待审 1已上墙(VIP 直接上墙)
}

// ListReq 评论列表(顶层分页, 每条带回复; 公开)。
type ListReq struct {
	g.Meta    `path:"/comment/list" method:"post" tags:"Front/Comment" summary:"评论列表"`
	MediaType int   `json:"media_type" v:"required|in:1,2,4,7#资源类型必填|资源类型非法"`
	ContentId int64 `json:"content_id" v:"required|min:1#内容ID必填"`
	Page      int   `json:"page"`
	Size      int   `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}
