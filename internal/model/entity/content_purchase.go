// Code maintained manually (内容购买记录, 付费解锁).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// ContentPurchase media_type: 1视频 2帖子 3漫画 4小说 5图集。
type ContentPurchase struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	MediaType int         `json:"mediaType" orm:"media_type"`
	ContentId int64       `json:"contentId" orm:"content_id"`
	Title     string      `json:"title"     orm:"title"`
	Amount    float64     `json:"amount"    orm:"amount"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
