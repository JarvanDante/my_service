// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserShareLog is the golang structure for table user_share_log.
type UserShareLog struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	Type      string      `json:"type"      orm:"type"       description:""` //
	TargetId  int64       `json:"targetId"  orm:"target_id"  description:""` //
	Channel   string      `json:"channel"   orm:"channel"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
