// Package v1 后台UGC投稿契约(列表 + 审核)。
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
	AuditBy      int64    `json:"audit_by"`
	AuditAt      string   `json:"audit_at"`
	CreatedAt    string   `json:"created_at"`
}

// ListReq 投稿列表。status/user_id/type 一律 string 接收, 空=不筛选。
type ListReq struct {
	g.Meta  `path:"/publishes" method:"get" tags:"Backend/Publish" summary:"投稿列表"`
	Status  string `json:"status"`
	UserId  string `json:"user_id"`
	Type    string `json:"type"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// AuditReq 审核: pass=true 通过, false 拒绝(必须给理由)。仅待审的投稿可审。
type AuditReq struct {
	g.Meta       `path:"/publishes/{id}/audit" method:"post" tags:"Backend/Publish" summary:"审核投稿"`
	Id           int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Pass         bool   `json:"pass"`
	RejectReason string `json:"reject_reason" v:"max-length:255#拒绝原因过长"`
}
type AuditRes struct{}
