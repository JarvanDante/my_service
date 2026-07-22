// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTask is the golang structure of table user_task for DAO operations like Where/Data.
type UserTask struct {
	g.Meta      `orm:"table:user_task, do:true"`
	Id          any         //
	Name        any         //
	Type        any         //
	Description any         //
	MaxNum      any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
