// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTask is the golang structure for table user_task.
type UserTask struct {
	Id          int64       `json:"id"          orm:"id"          description:""` //
	Name        string      `json:"name"        orm:"name"        description:""` //
	Type        string      `json:"type"        orm:"type"        description:""` //
	Description string      `json:"description" orm:"description" description:""` //
	MaxNum      int         `json:"maxNum"      orm:"max_num"     description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:""` //
}
