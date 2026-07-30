// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatConversation is the golang structure of table chat_conversation for DAO operations like Where/Data.
type ChatConversation struct {
	g.Meta      `orm:"table:chat_conversation, do:true"`
	Id          any         //
	SiteId      any         //
	UserId      any         //
	PeerId      any         //
	LastContent any         //
	LastAt      *gtime.Time //
	Unread      any         //
	Deleted     any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
