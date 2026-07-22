// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserSign is the golang structure of table user_sign for DAO operations like Where/Data.
type UserSign struct {
	g.Meta    `orm:"table:user_sign, do:true"`
	Id        any         //
	UserId    any         //
	YearMonth any         //
	Days      any         //
	Exchanges any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
