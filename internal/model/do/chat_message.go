// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatMessage is the golang structure of table chat_message for DAO operations like Where/Data.
type ChatMessage struct {
	g.Meta    `orm:"table:chat_message, do:true"`
	Id        any         //
	SiteId    any         //
	FromId    any         //
	ToId      any         //
	Content   any         //
	CreatedAt *gtime.Time //
}
