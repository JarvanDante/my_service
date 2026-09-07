package entity

import "github.com/gogf/gf/v2/os/gtime"

type PromoChannel struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Code      string      `json:"code"      orm:"code"`
	Name      string      `json:"name"      orm:"name"`
	Remark    string      `json:"remark"    orm:"remark"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
