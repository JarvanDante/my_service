// Code maintained manually (用户意见反馈).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Feedback struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	UserId      int64       `json:"userId"      orm:"user_id"`
	Type        int         `json:"type"        orm:"type"`
	ProblemType int         `json:"problemType" orm:"problem_type"`
	Content     string      `json:"content"     orm:"content"`
	Pics        string      `json:"pics"        orm:"pics"` // 原始 JSON 文本
	SysInfo     string      `json:"sysInfo"     orm:"sys_info"`
	MediaId     int64       `json:"mediaId"     orm:"media_id"`
	MediaTitle  string      `json:"mediaTitle"  orm:"media_title"`
	Status      int         `json:"status"      orm:"status"`
	Reply       string      `json:"reply"       orm:"reply"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
