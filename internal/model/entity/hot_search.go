// Code maintained manually (热搜词).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type HotSearch struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	Keyword     string      `json:"keyword"     orm:"keyword"`
	Heat        int         `json:"heat"        orm:"heat"`
	SearchCount int64       `json:"searchCount" orm:"search_count"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
