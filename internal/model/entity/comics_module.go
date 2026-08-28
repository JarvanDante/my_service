// Code maintained manually (漫画首页运营模块).
package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	ComicsModuleStyleHeroMix    = 1 // 1大2小 横图
	ComicsModuleStyleTwoWide    = 2 // 2小 横图
	ComicsModuleStyleOneWide    = 3 // 1大 横图
	ComicsModuleStyleTwoPoster  = 4 // 2竖图
	ComicsModuleStylePosterRail = 5 // 竖图横滑
	ComicsModuleStyleWideRail   = 6 // 横图横滑
	ComicsModuleStylePoster3x3  = 7 // 竖图3X3
)

const (
	ComicsModuleIconNew  = 1
	ComicsModuleIconStar = 2
	ComicsModuleIconFire = 3
)

const ComicsModulePosHome = "comic_home"

type ComicsModule struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Position  string      `json:"position"  orm:"position"`
	Style     int         `json:"style"     orm:"style"`
	Icon      int         `json:"icon"      orm:"icon"`
	TagIds    string      `json:"tagIds"    orm:"tag_ids"`
	Size      int         `json:"size"      orm:"size"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
