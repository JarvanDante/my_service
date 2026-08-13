// Package v1 前台UGC投稿契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64    `json:"id"`
	UserId       int64    `json:"user_id"`
	Type         int      `json:"type"` // 1视频 2漫画 3小说 4图集
	Title        string   `json:"title"`
	Intro        string   `json:"intro"`
	Cover        string   `json:"cover"`
	Resource     []string `json:"resource"`
	Tags         []string `json:"tags"`
	Status       int      `json:"status"` // 0待审 1通过 2拒绝 3已撤回
	RejectReason string   `json:"reject_reason"`
	AuditAt      string   `json:"audit_at"`
	CreatedAt    string   `json:"created_at"`
}

// SubmitReq 投稿(需登录)。标题/简介过敏感词, 落库即为待审(status=0)。
type SubmitReq struct {
	g.Meta   `path:"/publish/submit" method:"post" tags:"Front/Publish" summary:"提交投稿"`
	Type     int      `json:"type" v:"required|in:1,2,3,4#投稿类型必填|投稿类型非法"`
	Title    string   `json:"title" v:"required|max-length:128#标题必填|标题过长"`
	Intro    string   `json:"intro" v:"max-length:1000#简介过长"`
	Cover    string   `json:"cover"`
	Resource []string `json:"resource"` // 附件/外链
	Tags     []string `json:"tags"`
}
type SubmitRes struct {
	Id int64 `json:"id"`
}

// MyReq 我的投稿(需登录)。status 用 string 接收: 空=全部, 否则按值筛选。
type MyReq struct {
	g.Meta `path:"/publish/my" method:"get" tags:"Front/Publish" summary:"我的投稿"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type MyRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// CancelReq 撤回自己的待审投稿(已审核过的不可撤回)。
type CancelReq struct {
	g.Meta `path:"/publish/cancel" method:"post" tags:"Front/Publish" summary:"撤回投稿"`
	Id     int64 `json:"id" v:"required|min:1#投稿ID必填"`
}
type CancelRes struct{}
