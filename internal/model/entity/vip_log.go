// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// VipLog is the golang structure for table vip_log.
type VipLog struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	SiteId    int64       `json:"siteId"    orm:"site_id"    description:""` //
	UserId    int64       `json:"userId"    orm:"user_id"    description:""` //
	PackageId int64       `json:"packageId" orm:"package_id" description:""` //
	Days      int         `json:"days"      orm:"days"       description:""` //
	Price     float64     `json:"price"     orm:"price"      description:""` //
	StartAt   int64       `json:"startAt"   orm:"start_at"   description:""` //
	EndAt     int64       `json:"endAt"     orm:"end_at"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
