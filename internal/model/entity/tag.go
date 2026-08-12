// Code maintained manually (内容标签).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Tag struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	ContentType int         `json:"contentType" orm:"content_type"` // 1影片 2抖音 3动漫 4漫画 5图集 6帖子 7小说
	Name        string      `json:"name"        orm:"name"`
	Rank        int         `json:"rank"        orm:"rank"`
	Status      int         `json:"status"      orm:"status"` // 1启用 0禁用
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
