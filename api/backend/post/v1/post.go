// Package v1 后台帖子契约(审核/管理)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64    `json:"id"`
	UserId       int64    `json:"user_id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Pics         []string `json:"pics"`
	Topics       []string `json:"topics"`
	VideoUrl     string   `json:"video_url"`
	MediaId      int64    `json:"media_id"`
	ViewCount    int64    `json:"view_count"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	Status       int      `json:"status"`
	RejectReason string   `json:"reject_reason"`
	CreatedAt    string   `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/post" method:"get" tags:"Backend/Post" summary:"帖子列表"`
	Status  string `json:"status"`  // 空=全部  0待审 1通过 2拒绝 3用户删除
	Keyword string `json:"keyword"` // 标题模糊
	UserId  int64  `json:"user_id"` // 0=全部
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// AuditReq 审核(1通过 2拒绝需原因)。
type AuditReq struct {
	g.Meta `path:"/post/{id}/audit" method:"post" tags:"Backend/Post" summary:"审核帖子"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"` // 拒绝原因
}
type AuditRes struct{}

type DeleteReq struct {
	g.Meta `path:"/post/{id}" method:"delete" tags:"Backend/Post" summary:"删除帖子(硬删, 连带评论)"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
