// Package v1 后台数据统计接口契约(B8)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// 概览
type OverviewReq struct {
	g.Meta `path:"/stats/overview" method:"get" tags:"Backend/Stats" summary:"数据概览"`
}
type OverviewRes struct {
	TotalUsers      int     `json:"total_users"`       // 用户总数
	TodayNew        int     `json:"today_new"`         // 今日新增
	YesterdayNew    int     `json:"yesterday_new"`     // 昨日新增
	TodayActive     int     `json:"today_active"`      // 今日活跃(DAU)
	TodayPaidAmount float64 `json:"today_paid_amount"` // 今日充值金额
	TodayPaidOrders int     `json:"today_paid_orders"` // 今日充值单数
	YestPaidAmount  float64 `json:"yest_paid_amount"`  // 昨日充值金额
	TotalPaidAmount float64 `json:"total_paid_amount"` // 累计充值金额
	TotalPaidOrders int     `json:"total_paid_orders"` // 累计充值单数
}

// 注册趋势
type UserTrendItem struct {
	Date     int `json:"date"` // YYYYMMDD
	NewUsers int `json:"new_users"`
}
type UserTrendReq struct {
	g.Meta `path:"/stats/user-trend" method:"get" tags:"Backend/Stats" summary:"注册趋势"`
	Days   int `json:"days" v:"min:0|max:90#days 不合法|days 最大90"` // 默认7
}
type UserTrendRes struct {
	List []UserTrendItem `json:"list"`
}

// 充值趋势
type RechargeTrendItem struct {
	Date   int     `json:"date"`
	Orders int     `json:"orders"`
	Amount float64 `json:"amount"`
}
type RechargeTrendReq struct {
	g.Meta `path:"/stats/recharge-trend" method:"get" tags:"Backend/Stats" summary:"充值趋势"`
	Days   int `json:"days" v:"min:0|max:90#days 不合法|days 最大90"`
}
type RechargeTrendRes struct {
	List []RechargeTrendItem `json:"list"`
}

// 渠道分析
type ChannelStatItem struct {
	Channel       string  `json:"channel"`
	UserCount     int     `json:"user_count"`
	TotalRecharge float64 `json:"total_recharge"`
}
type ChannelStatsReq struct {
	g.Meta `path:"/stats/channels" method:"get" tags:"Backend/Stats" summary:"渠道分析"`
}
type ChannelStatsRes struct {
	List []ChannelStatItem `json:"list"`
}
