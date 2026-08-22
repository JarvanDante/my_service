package entity

import "github.com/gogf/gf/v2/os/gtime"

type OfficialGroup struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Intro     string      `json:"intro"     orm:"intro"`
	Avatar    string      `json:"avatar"    orm:"avatar"`
	Link      string      `json:"link"      orm:"link"`
	Platform  string      `json:"platform"  orm:"platform"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
