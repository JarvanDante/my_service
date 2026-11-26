// Code maintained manually (短剧分类运营模块).
package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	DramaModuleStyleHeroMix    = 1
	DramaModuleStyleTwoWide    = 2
	DramaModuleStyleOneWide    = 3
	DramaModuleStyleTwoPoster  = 4
	DramaModuleStylePosterRail = 5
	DramaModuleStyleWideRail   = 6
	DramaModuleStylePoster3x3  = 7
)

const (
	DramaModuleIconNew  = 1
	DramaModuleIconStar = 2
	DramaModuleIconFire = 3
)

const DramaModulePosHome = "drama_home"

type DramaModule struct {
	Id          int64       `json:"id"        orm:"id"`
	SiteId      int64       `json:"siteId"    orm:"site_id"`
	Name        string      `json:"name"      orm:"name"`
	Position    string      `json:"position"  orm:"position"`
	Style       int         `json:"style"     orm:"style"`
	Icon        int         `json:"icon"        orm:"icon"`
	CategoryIds string      `json:"categoryIds" orm:"category_ids"`
	TagIds      string      `json:"tagIds"      orm:"tag_ids"`
	Filter      string      `json:"filter"      orm:"filter"`
	Size        int         `json:"size"      orm:"size"`
	Rank        int         `json:"rank"      orm:"rank"`
	Status      int         `json:"status"    orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
