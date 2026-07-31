// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type SystemNotice struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Title     string      `json:"title"     orm:"title"`
	Content   string      `json:"content"   orm:"content"`
	Type      string      `json:"type"      orm:"type"`
	Status    int         `json:"status"    orm:"status"`
	CreatedBy int64       `json:"createdBy" orm:"created_by"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
