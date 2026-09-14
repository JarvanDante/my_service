// Code maintained manually (小说分类).
package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	NovelCategoryKindNormal = 0
	NovelCategoryKindNew    = 1
	NovelCategoryKindHot    = 2
	NovelCategoryKindRank   = 3
)

type NovelCategory struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Kind      int         `json:"kind"      orm:"kind"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
