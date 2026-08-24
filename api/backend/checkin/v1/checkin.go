package v1

import "github.com/gogf/gf/v2/frame/g"

type Config struct {
	MakeupPoints int    `json:"makeup_points"`
	MakeupLimit  int    `json:"makeup_limit"`
	MakeupDesc   string `json:"makeup_desc"`
	VipGroupId   int64  `json:"vip_group_id"`
}

type RewardRow struct {
	DayNum      int    `json:"day_num"`
	Label       string `json:"label"`
	Points      int64  `json:"points"`
	Gold        int64  `json:"gold"`
	VipDays     int    `json:"vip_days"`
	IsMilestone int    `json:"is_milestone"`
	MsPoints    int64  `json:"ms_points"`
	MsGold      int64  `json:"ms_gold"`
	MsVipDays   int    `json:"ms_vip_days"`
}

type GetConfigReq struct {
	g.Meta `path:"/checkin/config" method:"get" tags:"Backend/Checkin" summary:"签到全局配置"`
}
type GetConfigRes struct {
	Config Config `json:"config"`
}

type SaveConfigReq struct {
	g.Meta       `path:"/checkin/config" method:"put" tags:"Backend/Checkin" summary:"保存签到全局配置"`
	MakeupPoints int    `json:"makeup_points"`
	MakeupLimit  int    `json:"makeup_limit"`
	MakeupDesc   string `json:"makeup_desc"`
	VipGroupId   int64  `json:"vip_group_id"`
}
type SaveConfigRes struct{}

type RewardListReq struct {
	g.Meta `path:"/checkin/rewards" method:"get" tags:"Backend/Checkin" summary:"1-15天签到奖励"`
}
type RewardListRes struct {
	List []RewardRow `json:"list"`
}

type SaveRewardsReq struct {
	g.Meta `path:"/checkin/rewards" method:"put" tags:"Backend/Checkin" summary:"保存1-15天签到奖励"`
	List   []RewardRow `json:"list"`
}
type SaveRewardsRes struct{}
