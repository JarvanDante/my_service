// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCodeLog is the golang structure for table user_code_log.
type UserCodeLog struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	CodeId    int64       `json:"codeId"    orm:"code_id"    description:""` //
	Code      string      `json:"code"      orm:"code"       description:""` //
	CodeKey   string      `json:"codeKey"   orm:"code_key"   description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Type      string      `json:"type"      orm:"type"       description:""` //
	ObjectId  int64       `json:"objectId"  orm:"object_id"  description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	Username  string      `json:"username"  orm:"username"   description:""` //
	AddNum    int         `json:"addNum"    orm:"add_num"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
