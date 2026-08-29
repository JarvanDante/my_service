// Package v1 前台评论契约(移植自 tianbi comment)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id            int64    `json:"id"`
	UserId        int64    `json:"user_id"`
	Nickname      string   `json:"nickname"`
	Img           string   `json:"img"`
	IsVip         bool     `json:"is_vip"`
	ParentId      int64    `json:"parent_id"`
	RootId        int64    `json:"root_id"`
	ReplyUserId   int64    `json:"reply_user_id"`
	ReplyNickname string   `json:"reply_nickname"`
	Content       string   `json:"content"`
	Pics          []string `json:"pics"`
	LikeCount     int      `json:"like_count"`
	ReplyCount    int      `json:"reply_count"`
	Liked         bool     `json:"liked"`
	CreatedAt     string   `json:"created_at"`
	Replies       []Item   `json:"replies,omitempty"`
}

// AddReq 发表评论/回复(需登录, 过敏感词)。
type AddReq struct {
	g.Meta    `path:"/comment/add" method:"post" tags:"Front/Comment" summary:"发表评论"`
	MediaType int    `json:"media_type" v:"required|in:1,2,4,7#资源类型必填|资源类型非法"`
	ContentId int64  `json:"content_id" v:"required|min:1#内容ID必填"`
	ParentId  int64    `json:"parent_id"` // 0=顶层, >0=回复
	Content   string   `json:"content" v:"max-length:1000#内容过长"`
	Pics      []string `json:"pics"`
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
	Sort      int   `json:"sort"` // 0最新 1最热
	Page      int   `json:"page"`
	Size      int   `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// LikeReq 评论点赞/取消(需登录)。
type LikeReq struct {
	g.Meta `path:"/comment/like" method:"post" tags:"Front/Comment" summary:"评论点赞"`
	Id     int64 `json:"id" v:"required|min:1#评论ID必填"`
	Flag   bool  `json:"flag"`
}
type LikeRes struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}
