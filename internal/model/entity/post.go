// Code maintained manually (社区帖子).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Post struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	UserId       int64       `json:"userId"       orm:"user_id"`
	Title        string      `json:"title"        orm:"title"`
	Content      string      `json:"content"      orm:"content"`
	Pics         string      `json:"pics"         orm:"pics"` // 原始 JSON 数组文本
	Topics       string      `json:"topics"       orm:"topics"` // 话题名 JSON 数组
	Category     string      `json:"category"     orm:"category"`
	VideoUrl     string      `json:"videoUrl"     orm:"video_url"`
	MediaId      int64       `json:"mediaId"      orm:"media_id"`
	ViewCount    int64       `json:"viewCount"    orm:"view_count"`
	LikeCount    int         `json:"likeCount"    orm:"like_count"`
	CommentCount int         `json:"commentCount" orm:"comment_count"`
	Status       int         `json:"status"       orm:"status"` // 0待审 1通过 2拒绝 3删除
	RejectReason string      `json:"rejectReason" orm:"reject_reason"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}
