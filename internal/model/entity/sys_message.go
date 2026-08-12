// Code maintained manually (系统消息 + 已读记录).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type SysMessage struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"` // 0=全员
	Type      int         `json:"type"      orm:"type"`
	Title     string      `json:"title"     orm:"title"`
	Content   string      `json:"content"   orm:"content"`
	Status    int         `json:"status"    orm:"status"` // 1发布 0撤回
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

type SysMessageRead struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	MessageId int64       `json:"messageId" orm:"message_id"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
