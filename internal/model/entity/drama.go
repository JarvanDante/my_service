// Code maintained manually (短剧 + 剧集 + 分类).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// 剧集付费方式：0 跟随整部默认，1 免费，2 金币，3 VIP。
const (
	DramaPayFollow = 0
	DramaPayFree   = 1
	DramaPayCoin   = 2
	DramaPayVip    = 3
)

type Drama struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	Title        string      `json:"title"        orm:"title"`
	Cover        string      `json:"cover"        orm:"cover"`
	Intro        string      `json:"intro"        orm:"intro"`
	Category     string      `json:"category"     orm:"category"`
	Tags         string      `json:"tags"         orm:"tags"` // jsonb 原文
	IsVip        int         `json:"isVip"        orm:"is_vip"`
	Price        float64     `json:"price"        orm:"price"`
	FreeEpisode  int         `json:"freeEpisode"  orm:"free_episode"`
	EpisodeCount int         `json:"episodeCount" orm:"episode_count"`
	ViewCount    int64       `json:"viewCount"    orm:"view_count"`
	BuyCount     int64       `json:"buyCount"     orm:"buy_count"`
	LikeCount    int64       `json:"likeCount"    orm:"like_count"`
	UpdateStatus int         `json:"updateStatus" orm:"update_status"`
	Rank         int         `json:"rank"         orm:"rank"`
	IsRecommend  int         `json:"isRecommend"  orm:"is_recommend"`
	Status       int         `json:"status"       orm:"status"`
	MediaCode    string      `json:"mediaCode"    orm:"media_code"`
	PublishedAt  *gtime.Time `json:"publishedAt"  orm:"published_at"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}

type DramaEpisode struct {
	Id         int64       `json:"id"         orm:"id"`
	SiteId     int64       `json:"siteId"     orm:"site_id"`
	DramaId    int64       `json:"dramaId"    orm:"drama_id"`
	Seq        int         `json:"seq"        orm:"seq"`
	Title      string      `json:"title"      orm:"title"`
	MediaCode  string      `json:"mediaCode"  orm:"media_code"`
	Duration   int         `json:"duration"   orm:"duration"`
	PayMode    int         `json:"payMode"    orm:"pay_mode"`
	Price      float64     `json:"price"      orm:"price"`
	Status     int         `json:"status"     orm:"status"`
	MediaReady int         `json:"mediaReady" orm:"media_ready"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}

const (
	DramaCategoryKindNormal = 0
	DramaCategoryKindNew    = 1
	DramaCategoryKindHot    = 2
	DramaCategoryKindRank   = 3
)

type DramaCategory struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Kind      int         `json:"kind"      orm:"kind"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
