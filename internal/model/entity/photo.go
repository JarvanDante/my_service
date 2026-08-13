// Code maintained manually (图集).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// PhotoAlbum 图集。与 Comics 不同, 图片直接挂主表 pics(jsonb), 没有章节从表;
// 上下架状态复用 ContentStatus* 那套编码(见 comics.go)。
type PhotoAlbum struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Title     string      `json:"title"     orm:"title"`
	Cover     string      `json:"cover"     orm:"cover"`
	Intro     string      `json:"intro"     orm:"intro"`
	Category  string      `json:"category"  orm:"category"`
	Tags      string      `json:"tags"      orm:"tags"` // jsonb 原文
	IsVip     int         `json:"isVip"     orm:"is_vip"`
	Price     float64     `json:"price"     orm:"price"`
	FreeCount int         `json:"freeCount" orm:"free_count"` // 未解锁可预览的前 N 张
	Pics      string      `json:"pics"      orm:"pics"`       // jsonb 原文
	PicCount  int         `json:"picCount"  orm:"pic_count"`
	ViewCount int64       `json:"viewCount" orm:"view_count"`
	BuyCount  int64       `json:"buyCount"  orm:"buy_count"`
	LikeCount int64       `json:"likeCount" orm:"like_count"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	PublishId int64       `json:"publishId" orm:"publish_id"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
