// Code maintained manually (推广应用).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Application struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	Name        string      `json:"name"        orm:"name"`
	Tag         int         `json:"tag"         orm:"tag"`
	Intro       string      `json:"intro"       orm:"intro"`
	Avatar      string      `json:"avatar"      orm:"avatar"`
	DownloadUrl string      `json:"downloadUrl" orm:"download_url"`
	IosUrl      string      `json:"iosUrl"      orm:"ios_url"`
	AndroidUrl  string      `json:"androidUrl"  orm:"android_url"`
	LocIds      string      `json:"locIds"      orm:"loc_ids"` // 原始 JSON 数组文本
	Rank        int         `json:"rank"        orm:"rank"`
	DownTotal   int64       `json:"downTotal"   orm:"down_total"`
	Status      int         `json:"status"      orm:"status"` // 1上架 0下架
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
