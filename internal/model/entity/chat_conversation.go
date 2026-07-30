// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatConversation is the golang structure for table chat_conversation.
type ChatConversation struct {
	Id          int64       `json:"id"          orm:"id"           description:""` //
	SiteId      int64       `json:"siteId"      orm:"site_id"      description:""` //
	UserId      int64       `json:"userId"      orm:"user_id"      description:""` //
	PeerId      int64       `json:"peerId"      orm:"peer_id"      description:""` //
	LastContent string      `json:"lastContent" orm:"last_content" description:""` //
	LastAt      *gtime.Time `json:"lastAt"      orm:"last_at"      description:""` //
	Unread      int         `json:"unread"      orm:"unread"       description:""` //
	Deleted     int         `json:"deleted"     orm:"deleted"      description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""` //
}
