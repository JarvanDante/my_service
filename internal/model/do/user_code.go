// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCode is the golang structure of table user_code for DAO operations like Where/Data.
type UserCode struct {
	g.Meta    `orm:"table:user_code, do:true"`
	Id        any         //
	Name      any         //
	Code      any         //
	CodeKey   any         //
	Type      any         //
	ObjectId  any         //
	AddNum    any         //
	CanUseNum any         //
	UsedNum   any         //
	Status    any         //
	ExpiredAt any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
