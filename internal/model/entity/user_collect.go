// Code maintained manually (用户收藏/点赞/踩).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type UserCollect struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	OpType    int         `json:"opType"    orm:"op_type"` // 1收藏 2点赞 3踩
	MediaType int         `json:"mediaType" orm:"media_type"`
	ContentId int64       `json:"contentId" orm:"content_id"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
