// Code maintained manually (内容评论).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Comment struct {
	Id         int64       `json:"id"         orm:"id"`
	SiteId     int64       `json:"siteId"     orm:"site_id"`
	UserId     int64       `json:"userId"     orm:"user_id"`
	MediaType  int         `json:"mediaType"  orm:"media_type"`
	ContentId  int64       `json:"contentId"  orm:"content_id"`
	ParentId   int64       `json:"parentId"   orm:"parent_id"`
	RootId     int64       `json:"rootId"     orm:"root_id"`
	Content    string      `json:"content"    orm:"content"`
	LikeCount  int         `json:"likeCount"  orm:"like_count"`
	ReplyCount int         `json:"replyCount" orm:"reply_count"`
	Status     int         `json:"status"     orm:"status"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
}
