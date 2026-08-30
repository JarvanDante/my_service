package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	KingkongPosComics  = "comics"
	KingkongPosCartoon = "cartoon"
	KingkongPosMovie   = "movie"
	KingkongPosNovel   = "novel"
	KingkongPosShort   = "short"

	KingkongModeBlock  = "block"
	KingkongModeList   = "list"
	KingkongModeDouyin = "douyin"
)

type KingkongItem struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	IconUrl   string      `json:"iconUrl"   orm:"icon_url"`
	OpenMode  string      `json:"openMode"  orm:"open_mode"`
	Link      string      `json:"link"      orm:"link"`
	AppLink   string      `json:"appLink"   orm:"app_link"`
	Position  string      `json:"position"  orm:"position"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
