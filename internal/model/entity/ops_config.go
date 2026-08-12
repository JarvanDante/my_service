// Code maintained manually (运营配置: 公告/跳转位/敏感词).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Announcement struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Title     string      `json:"title"     orm:"title"`
	Content   string      `json:"content"   orm:"content"`
	TextNode  string      `json:"textNode"  orm:"text_node"`
	Cover     string      `json:"cover"     orm:"cover"`
	JumpUrl   string      `json:"jumpUrl"   orm:"jump_url"`
	SysType   string      `json:"sysType"   orm:"sys_type"`
	StartAt   *gtime.Time `json:"startAt"   orm:"start_at"`
	EndAt     *gtime.Time `json:"endAt"     orm:"end_at"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

type Jumptab struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	CnName      string      `json:"cnName"      orm:"cn_name"`
	EnName      string      `json:"enName"      orm:"en_name"`
	Avatar      string      `json:"avatar"      orm:"avatar"`
	Link        string      `json:"link"        orm:"link"`
	PicJumpLink string      `json:"picJumpLink" orm:"pic_jump_link"`
	Location    int         `json:"location"    orm:"location"`
	Rank        int         `json:"rank"        orm:"rank"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

type FilterWord struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Word      string      `json:"word"      orm:"word"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
