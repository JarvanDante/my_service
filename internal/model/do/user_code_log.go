// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCodeLog is the golang structure of table user_code_log for DAO operations like Where/Data.
type UserCodeLog struct {
	g.Meta    `orm:"table:user_code_log, do:true"`
	Id        any         //
	CodeId    any         //
	Code      any         //
	CodeKey   any         //
	Name      any         //
	Type      any         //
	ObjectId  any         //
	UserId    any         //
	Username  any         //
	AddNum    any         //
	CreatedAt *gtime.Time //
}
