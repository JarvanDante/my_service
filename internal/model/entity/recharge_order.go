// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RechargeOrder is the golang structure for table recharge_order.
type RechargeOrder struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	OrderNo   string      `json:"orderNo"   orm:"order_no"   description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	PackageId int64       `json:"packageId" orm:"package_id" description:""` //
	Amount    float64     `json:"amount"    orm:"amount"     description:""` //
	Coin      float64     `json:"coin"      orm:"coin"       description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	PayAt     *gtime.Time `json:"payAt"     orm:"pay_at"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
