// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type AdminUser struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	Username    string      `json:"username"    orm:"username"`
	Password    string      `json:"password"    orm:"password"`
	Salt        string      `json:"salt"        orm:"salt"`
	Nickname    string      `json:"nickname"    orm:"nickname"`
	RoleId      int64       `json:"roleId"      orm:"role_id"`
	Status      int         `json:"status"      orm:"status"`
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at"`
	LastIp      string      `json:"lastIp"      orm:"last_ip"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
