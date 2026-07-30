// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatMessage is the golang structure for table chat_message.
type ChatMessage struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	FromId    int64       `json:"fromId"    orm:"from_id"    description:""` //
	ToId      int64       `json:"toId"      orm:"to_id"      description:""` //
	Content   string      `json:"content"   orm:"content"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
