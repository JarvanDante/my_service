// Code maintained manually (漫画作品 + 章节).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// 作品上下架状态(漫画/小说/图集共用同一套编码)
const (
	ContentStatusPending = 0 // 待上架/待审核
	ContentStatusOnline  = 1 // 上架
	ContentStatusOffline = 2 // 下架
)

type Comics struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	Title        string      `json:"title"        orm:"title"`
	Author       string      `json:"author"       orm:"author"`
	Cover        string      `json:"cover"        orm:"cover"`
	Intro        string      `json:"intro"        orm:"intro"`
	Category     string      `json:"category"     orm:"category"`
	Tags         string      `json:"tags"         orm:"tags"` // jsonb 原文
	IsVip        int         `json:"isVip"        orm:"is_vip"`
	Price        float64     `json:"price"        orm:"price"`
	FreeChapter  int         `json:"freeChapter"  orm:"free_chapter"`
	ChapterCount int         `json:"chapterCount" orm:"chapter_count"`
	ViewCount    int64       `json:"viewCount"    orm:"view_count"`
	BuyCount     int64       `json:"buyCount"     orm:"buy_count"`
	LikeCount    int64       `json:"likeCount"    orm:"like_count"`
	UpdateStatus int         `json:"updateStatus" orm:"update_status"`
	Rank         int         `json:"rank"         orm:"rank"`
	IsRecommend  int         `json:"isRecommend"  orm:"is_recommend"`
	Status       int         `json:"status"       orm:"status"`
	PublishId    int64       `json:"publishId"    orm:"publish_id"`
	MediaCode    string      `json:"mediaCode"    orm:"media_code"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}

type ComicsChapter struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	ComicsId  int64       `json:"comicsId"  orm:"comics_id"`
	Seq       int         `json:"seq"       orm:"seq"`
	Title     string      `json:"title"     orm:"title"`
	Pics      string      `json:"pics"      orm:"pics"` // jsonb 原文
	PicCount  int         `json:"picCount"  orm:"pic_count"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
