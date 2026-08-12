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

// UnreadReq 未读数(需登录)。
type UnreadReq struct {
	g.Meta `path:"/message/unread" method:"get" tags:"Front/Message" summary:"未读消息数"`
}
type UnreadRes struct {
	Count int `json:"count"`
}

// ReadReq 标记已读(id>0 单条, all=true 全部, 需登录)。
type ReadReq struct {
	g.Meta `path:"/message/read" method:"post" tags:"Front/Message" summary:"标记已读"`
	Id     int64 `json:"id"`
	All    bool  `json:"all"`
}
type ReadRes struct{}
