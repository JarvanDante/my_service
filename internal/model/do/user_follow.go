// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFollow is the golang structure of table user_follow for DAO operations like Where/Data.
type UserFollow struct {
	g.Meta    `orm:"table:user_follow, do:true"`
	Id        any         //
	UserId    any         //
	HomeId    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
