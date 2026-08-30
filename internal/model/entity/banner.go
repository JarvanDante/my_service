package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	BannerPosComics  = "comics"
	BannerPosCartoon = "cartoon"
	BannerPosVideo   = "video"
	BannerPosNovel   = "novel"
)

type Banner struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Position  string      `json:"position"  orm:"position"`
	Title     string      `json:"title"     orm:"title"`
	CoverUrl  string      `json:"coverUrl"  orm:"cover_url"`
	Link      string      `json:"link"      orm:"link"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
