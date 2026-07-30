// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type VipLog struct {
	Id        int64       `json:"id"        orm:"id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	PackageId int64       `json:"packageId" orm:"package_id"`
	Days      int         `json:"days"      orm:"days"`
	Price     float64     `json:"price"     orm:"price"`
	StartAt   int64       `json:"startAt"   orm:"start_at"`
	EndAt     int64       `json:"endAt"     orm:"end_at"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
