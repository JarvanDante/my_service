// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// VipPackage is the golang structure for table vip_package.
type VipPackage struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Days      int         `json:"days"      orm:"days"       description:""` //
	Price     float64     `json:"price"     orm:"price"      description:""` //
	GroupId   int64       `json:"groupId"   orm:"group_id"   description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
