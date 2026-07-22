// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTaskLog is the golang structure for table user_task_log.
type UserTaskLog struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	TaskId    int64       `json:"taskId"    orm:"task_id"    description:""` //
	Type      string      `json:"type"      orm:"type"       description:""` //
	Num       int         `json:"num"       orm:"num"        description:""` //
	LogDate   int         `json:"logDate"   orm:"log_date"   description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
