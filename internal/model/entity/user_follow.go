// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFollow is the golang structure for table user_follow.
type UserFollow struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	HomeId    int64       `json:"homeId"    orm:"home_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
