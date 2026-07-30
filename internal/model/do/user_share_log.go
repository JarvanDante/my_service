// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserShareLog is the golang structure of table user_share_log for DAO operations like Where/Data.
type UserShareLog struct {
	g.Meta    `orm:"table:user_share_log, do:true"`
	Id        any         //
	SiteId    any         //
	UserId    any         //
	Type      any         //
	TargetId  any         //
	Channel   any         //
	CreatedAt *gtime.Time //
}
