// Code maintained manually (签到阶梯奖励配置).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type CheckinReward struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	DayNum    int         `json:"dayNum"    orm:"day_num"`
	UserType  int         `json:"userType"  orm:"user_type"`
	Gold      int64       `json:"gold"      orm:"gold"`
	VipDays   int         `json:"vipDays"   orm:"vip_days"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
