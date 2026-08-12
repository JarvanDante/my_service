// Code maintained manually (应用基础配置 KV).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type AppConfig struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Grp       string      `json:"grp"       orm:"grp"` // 分组
	Key       string      `json:"key"       orm:"key"`
	Value     string      `json:"value"     orm:"value"` // 原始 JSON 文本
	Remark    string      `json:"remark"    orm:"remark"`
	Status    int         `json:"status"    orm:"status"` // 1启用 0禁用
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
