// Code maintained manually (用户UGC投稿).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// 投稿类型(与 paywall 的 media_type 是两套编码, 不要混用)
const (
	PublishTypeVideo  = 1
	PublishTypeComics = 2
	PublishTypeNovel  = 3
	PublishTypePhoto  = 4
)

// 投稿审核状态
const (
	PublishStatusPending  = 0 // 待审
	PublishStatusPass     = 1 // 通过
	PublishStatusReject   = 2 // 拒绝
	PublishStatusCanceled = 3 // 用户撤回(不删行, 留痕)
)

type UserPublish struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	UserId       int64       `json:"userId"       orm:"user_id"`
	Type         int         `json:"type"         orm:"type"`
	Title        string      `json:"title"        orm:"title"`
	Intro        string      `json:"intro"        orm:"intro"`
	Cover        string      `json:"cover"        orm:"cover"`
	Resource     string      `json:"resource"     orm:"resource"` // jsonb 原文
	Tags         string      `json:"tags"         orm:"tags"`     // jsonb 原文
	Status       int         `json:"status"       orm:"status"`
	RejectReason string      `json:"rejectReason" orm:"reject_reason"`
	AuditBy      int64       `json:"auditBy"      orm:"audit_by"`
	AuditAt      *gtime.Time `json:"auditAt"      orm:"audit_at"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}
