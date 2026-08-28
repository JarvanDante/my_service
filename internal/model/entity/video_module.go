// Code maintained manually (视频/动漫首页运营模块, 表结构相同)。
package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	VideoModulePosHome   = "video_home"
	CartoonModulePosHome = "cartoon_home"
)

type VideoModule struct {
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
