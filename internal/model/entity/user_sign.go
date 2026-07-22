// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserSign is the golang structure for table user_sign.
type UserSign struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	YearMonth int         `json:"yearMonth" orm:"year_month" description:""` //
	Days      string      `json:"days"      orm:"days"       description:""` //
	Exchanges string      `json:"exchanges" orm:"exchanges"  description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
