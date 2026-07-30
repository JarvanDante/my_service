// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RechargePackage is the golang structure for table recharge_package.
type RechargePackage struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Amount    float64     `json:"amount"    orm:"amount"     description:""` //
	Coin      float64     `json:"coin"      orm:"coin"       description:""` //
	Bonus     float64     `json:"bonus"     orm:"bonus"      description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
