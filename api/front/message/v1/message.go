// Package v1 前台系统消息契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type MsgItem struct {
	Id        int64  `json:"id"`
	Type      int    `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

// ListReq 我的消息列表(全员消息 + 发给我的, 需登录)。
type ListReq struct {
	g.Meta `path:"/message/list" method:"get" tags:"Front/Message" summary:"我的消息列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type ListRes struct {
	List  []MsgItem `json:"list"`
	Total int       `json:"total"`
}

// UnreadReq 未读数(需登录)。count=站内信+评论+点赞。
type UnreadReq struct {
	g.Meta `path:"/message/unread" method:"get" tags:"Front/Message" summary:"未读消息数"`
}
type UnreadRes struct {
	Count   int `json:"count"`
	Sys     int `json:"sys"`
	Comment int `json:"comment"`
	Like    int `json:"like"`
}

type InteractItem struct {
	Id            int64  `json:"id"`
	Channel       string `json:"channel"`
	SubType       string `json:"sub_type"`
	IsRead        bool   `json:"is_read"`
	CreatedAt     string `json:"created_at"`
	ActorId       int64  `json:"actor_id"`
	ActorName     string `json:"actor_name"`
	ActorAvatar   string `json:"actor_avatar"`
	ActorSex      int    `json:"actor_sex"`
	ActorIsVip    bool   `json:"actor_is_vip"`
	ActorCount    int    `json:"actor_count"`
	MediaType     int    `json:"media_type"`
	ContentId     int64  `json:"content_id"`
	ObjectTitle   string `json:"object_title"`
	TargetType    string `json:"target_type"`
	CommentId     int64  `json:"comment_id"`
	RootCommentId int64  `json:"root_comment_id"`
	Snippet       string `json:"snippet"`
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
	Deleted       bool   `json:"deleted"`
}

// InteractListReq 评论/点赞互动列表。channel=comment|like。
type InteractListReq struct {
	g.Meta  `path:"/message/interact" method:"get" tags:"Front/Message" summary:"互动消息列表"`
	Channel string `json:"channel"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type InteractListRes struct {
	List  []InteractItem `json:"list"`
	Total int            `json:"total"`
}

// InteractReadReq 互动已读。id>0 单条; all=true 按 channel 或全部互动。
type InteractReadReq struct {
	g.Meta  `path:"/message/interact/read" method:"post" tags:"Front/Message" summary:"互动消息已读"`
	Id      int64  `json:"id"`
	All     bool   `json:"all"`
	Channel string `json:"channel"`
}
type InteractReadRes struct{}

// ReadReq 标记已读(id>0 单条, all=true 全部, 需登录)。
type ReadReq struct {
	g.Meta `path:"/message/read" method:"post" tags:"Front/Message" summary:"标记已读"`
	Id     int64 `json:"id"`
	All    bool  `json:"all"`
}
type ReadRes struct{}
