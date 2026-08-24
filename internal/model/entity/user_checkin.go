// Code maintained manually (用户每日签到记录).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type UserCheckin struct {
	Id               int64       `json:"id"               orm:"id"`
	SiteId           int64       `json:"siteId"           orm:"site_id"`
	UserId           int64       `json:"userId"           orm:"user_id"`
	CheckinDate      *gtime.Time `json:"checkinDate"      orm:"checkin_date"`
	ContinuouslyDays int         `json:"continuouslyDays" orm:"continuously_days"`
	RewardGold       int64       `json:"rewardGold"       orm:"reward_gold"`
	RewardPoints     int64       `json:"rewardPoints"     orm:"reward_points"`
	RewardVipDays    int         `json:"rewardVipDays"    orm:"reward_vip_days"`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"`
}
