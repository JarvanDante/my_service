// Code maintained manually (漫画分类).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// 漫画分类栏位类型(H5 顶栏)
const (
	ComicsCategoryKindNormal = 0 // 普通分类, 按名称筛作品
	ComicsCategoryKindNew    = 1 // 新更
	ComicsCategoryKindHot    = 2 // 推荐
	ComicsCategoryKindRank   = 3 // 榜单
)

type ComicsCategory struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Kind      int         `json:"kind"      orm:"kind"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
