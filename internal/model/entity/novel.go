// Code maintained manually (小说作品 + 章节).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// 上下架状态复用 comics.go 里的 ContentStatus* 常量, 全站一套编码。

type Novel struct {
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
	WordCount    int64       `json:"wordCount"    orm:"word_count"` // 全书总字数, 由章节汇总
	IsAudio      int         `json:"isAudio"      orm:"is_audio"`
	ViewCount    int64       `json:"viewCount"    orm:"view_count"`
	BuyCount     int64       `json:"buyCount"     orm:"buy_count"`
	LikeCount    int64       `json:"likeCount"    orm:"like_count"`
	UpdateStatus int         `json:"updateStatus" orm:"update_status"`
	Rank         int         `json:"rank"         orm:"rank"`
	Status       int         `json:"status"       orm:"status"`
	PublishId    int64       `json:"publishId"    orm:"publish_id"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}

type NovelChapter struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	NovelId   int64       `json:"novelId"   orm:"novel_id"`
	Seq       int         `json:"seq"       orm:"seq"`
	Title     string      `json:"title"     orm:"title"`
	Content   string      `json:"content"   orm:"content"` // 正文全文
	WordCount int         `json:"wordCount" orm:"word_count"`
	AudioUrl  string      `json:"audioUrl"  orm:"audio_url"` // 无声小说为空
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
