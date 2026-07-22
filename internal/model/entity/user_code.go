// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCode is the golang structure for table user_code.
type UserCode struct {
	Id        int64       `json:"id"        orm:"id"          description:""` //
	Name      string      `json:"name"      orm:"name"        description:""` //
	Code      string      `json:"code"      orm:"code"        description:""` //
	CodeKey   string      `json:"codeKey"   orm:"code_key"    description:""` //
	Type      string      `json:"type"      orm:"type"        description:""` //
	ObjectId  int64       `json:"objectId"  orm:"object_id"   description:""` //
	AddNum    int         `json:"addNum"    orm:"add_num"     description:""` //
	CanUseNum int         `json:"canUseNum" orm:"can_use_num" description:""` //
	UsedNum   int         `json:"usedNum"   orm:"used_num"    description:""` //
	Status    int         `json:"status"    orm:"status"      description:""` //
	ExpiredAt int64       `json:"expiredAt" orm:"expired_at"  description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"  description:""` //
}
