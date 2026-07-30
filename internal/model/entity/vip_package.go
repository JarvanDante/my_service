// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type VipPackage struct {
	Id        int64       `json:"id"        orm:"id"`
	Name      string      `json:"name"      orm:"name"`
	Days      int         `json:"days"      orm:"days"`
	Price     float64     `json:"price"     orm:"price"`
	GroupId   int64       `json:"groupId"   orm:"group_id"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
