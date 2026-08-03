// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type AdminRole struct {
	Id          int64       `json:"id"        orm:"id"`
	SiteId      int64       `json:"siteId"    orm:"site_id"`
	Name        string      `json:"name"      orm:"name"`
	Code        string      `json:"code"      orm:"code"`
	Remark      string      `json:"remark"    orm:"remark"`
	Status      int         `json:"status"      orm:"status"`
	Permissions string      `json:"permissions" orm:"permissions"`
	CreatedAt   *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
