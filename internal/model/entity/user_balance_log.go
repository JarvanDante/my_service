// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserBalanceLog is the golang structure for table user_balance_log.
type UserBalanceLog struct {
	Id            int64       `json:"id"            orm:"id"             description:""` //
	UserId        int64       `json:"userId"        orm:"user_id"        description:""` //
	Direction     int         `json:"direction"     orm:"direction"      description:""` //
	Scene         string      `json:"scene"         orm:"scene"          description:""` //
	Amount        float64     `json:"amount"        orm:"amount"         description:""` //
	BalanceBefore float64     `json:"balanceBefore" orm:"balance_before" description:""` //
	BalanceAfter  float64     `json:"balanceAfter"  orm:"balance_after"  description:""` //
	RefId         string      `json:"refId"         orm:"ref_id"         description:""` //
	Remark        string      `json:"remark"        orm:"remark"         description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
}
