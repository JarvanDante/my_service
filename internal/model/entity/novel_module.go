// Code maintained manually (小说分类运营模块).
package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	NovelModuleStyleHeroMix    = 1
	NovelModuleStyleTwoWide    = 2
	NovelModuleStyleOneWide    = 3
	NovelModuleStyleTwoPoster  = 4
	NovelModuleStylePosterRail = 5
	NovelModuleStyleWideRail   = 6
	NovelModuleStylePoster3x3  = 7
)

const (
	NovelModuleIconNew  = 1
	NovelModuleIconStar = 2
	NovelModuleIconFire = 3
)

const NovelModulePosHome = "novel_home"

type NovelModule struct {
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
