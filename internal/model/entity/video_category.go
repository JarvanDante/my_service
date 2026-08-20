package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	VideoCategoryKindNormal = 0
	VideoCategoryKindNew    = 1
	VideoCategoryKindHot    = 2
	VideoCategoryKindRank   = 3
)

type VideoCategory struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Kind      int         `json:"kind"      orm:"kind"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
