// Package v1 前台签到接口契约(移植自 tianbi checkin)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type RewardItem struct {
	Gold    int64 `json:"gold"`
	VipDays int   `json:"vip_days"`
}
type RewardCfg struct {
	DayNum   int   `json:"day_num"`
	UserType int   `json:"user_type"`
	Gold     int64 `json:"gold"`
	VipDays  int   `json:"vip_days"`
}
type RecordItem struct {
	Date             string `json:"date"`
	ContinuouslyDays int    `json:"continuously_days"`
	RewardGold       int64  `json:"reward_gold"`
}

// Click 用户签到。
type ClickReq struct {
	g.Meta `path:"/checkin/click" method:"post" tags:"Front/Checkin" summary:"用户签到"`
}
type ClickRes struct {
	Message          string       `json:"message"`
	TodayChecked     bool         `json:"today_checked"`
	ContinuouslyDays int          `json:"continuously_days"`
	Rewards          []RewardItem `json:"rewards"`
}

// Info 签到规则(阶梯配置)+ 连续天数 + 记录。
type InfoReq struct {
	g.Meta `path:"/checkin/prize" method:"post" tags:"Front/Checkin" summary:"签到规则与记录"`
}
type InfoRes struct {
	Rewards          []RewardCfg  `json:"rewards"`
	TodayChecked     bool         `json:"today_checked"`
	ContinuouslyDays int          `json:"continuously_days"`
	Records          []RecordItem `json:"records"`
}
