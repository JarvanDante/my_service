package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	VideoStatusDraft     = 0
	VideoStatusPublished = 1
	VideoStatusOffline   = 2
)

const (
	VideoKindVideo   = 0
	VideoKindCartoon = 2
	VideoKindDouyin  = 3
)

type Video struct {
	Id            int64       `json:"id"            orm:"id"`
	SiteId        int64       `json:"siteId"        orm:"site_id"`
	Title         string      `json:"title"         orm:"title"`
	Description   string      `json:"description"   orm:"description"`
	CoverUrl      string      `json:"coverUrl"      orm:"cover_url"`
	CoverKey      string      `json:"coverKey"      orm:"cover_key"`
	CoverMediaId  int64       `json:"coverMediaId"  orm:"cover_media_id"`
	SourceUrl     string      `json:"sourceUrl"     orm:"source_url"`
	SourceKey     string      `json:"sourceKey"     orm:"source_key"`
	SourceMediaId int64       `json:"sourceMediaId" orm:"source_media_id"`
	MediaCode     string      `json:"mediaCode"     orm:"media_code"`
	Kind          int         `json:"kind"          orm:"kind"` // 0视频 2动漫 3抖音
	Category      string      `json:"category"      orm:"category"`
	Tags          string      `json:"tags"          orm:"tags"` // jsonb 原文
	Duration      int         `json:"duration"      orm:"duration"`
	Sort          int         `json:"sort"          orm:"sort"`
	Status        int         `json:"status"        orm:"status"`
	CreatedBy     int64       `json:"createdBy"     orm:"created_by"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"`
}
