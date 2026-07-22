// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTaskLog is the golang structure of table user_task_log for DAO operations like Where/Data.
type UserTaskLog struct {
	g.Meta    `orm:"table:user_task_log, do:true"`
	Id        any         //
	UserId    any         //
	TaskId    any         //
	Type      any         //
	Num       any         //
	LogDate   any         //
	CreatedAt *gtime.Time //
}
