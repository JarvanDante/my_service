// Code maintained manually (签到阶梯奖励配置).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type CheckinReward struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	DayNum      int         `json:"dayNum"      orm:"day_num"`
	UserType    int         `json:"userType"    orm:"user_type"`
	Label       string      `json:"label"       orm:"label"`
	Gold        int64       `json:"gold"        orm:"gold"`
	Points      int64       `json:"points"      orm:"points"`
	VipDays     int         `json:"vipDays"     orm:"vip_days"`
	IsMilestone int         `json:"isMilestone" orm:"is_milestone"`
	MsPoints    int64       `json:"msPoints"    orm:"ms_points"`
	MsGold      int64       `json:"msGold"      orm:"ms_gold"`
	MsVipDays   int         `json:"msVipDays"   orm:"ms_vip_days"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

type CheckinConfig struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	MakeupPoints int         `json:"makeupPoints" orm:"makeup_points"`
	MakeupLimit  int         `json:"makeupLimit"  orm:"makeup_limit"`
	MakeupDesc   string      `json:"makeupDesc"   orm:"makeup_desc"`
	VipGroupId   int64       `json:"vipGroupId"   orm:"vip_group_id"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}
