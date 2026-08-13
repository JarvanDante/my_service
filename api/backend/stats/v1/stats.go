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

// ---------- 扩展维度(分析页图表化) ----------

// 时段分布: 近 N 天内, 注册与支付订单按小时(0~23)聚合
type HourDistItem struct {
	Hour      int `json:"hour"`
	Registers int `json:"registers"` // 注册人数
	Orders    int `json:"orders"`    // 支付订单数
}
type HourDistReq struct {
	g.Meta `path:"/stats/hour-dist" method:"get" tags:"Backend/Stats" summary:"时段分布(注册/支付)"`
	Days   int `json:"days" v:"min:0|max:90#days 不合法|days 最大90"`
}
type HourDistRes struct {
	List []HourDistItem `json:"list"`
}

// 设备分布
type DeviceStatItem struct {
	DeviceType string `json:"device_type"`
	Count      int    `json:"count"`
}
type DeviceStatsReq struct {
	g.Meta `path:"/stats/devices" method:"get" tags:"Backend/Stats" summary:"设备分布"`
}
type DeviceStatsRes struct {
	List []DeviceStatItem `json:"list"`
}

// 内容库概览: 各内容类型的作品状态分布 + 消费情况
type ContentStatItem struct {
	MediaType int     `json:"media_type"` // 1视频 2帖子 3漫画 4小说 5图集
	TypeName  string  `json:"type_name"`
	Online    int     `json:"online"`     // 已上架/已通过
	Pending   int     `json:"pending"`    // 待上架/待审
	Offline   int     `json:"offline"`    // 下架/拒绝
	Views     int64   `json:"views"`      // 累计观看
	Buys      int     `json:"buys"`       // 累计购买件数(content_purchase)
	BuyAmount float64 `json:"buy_amount"` // 累计购买金额
}
type ContentStatsReq struct {
	g.Meta `path:"/stats/content" method:"get" tags:"Backend/Stats" summary:"内容库概览"`
}
type ContentStatsRes struct {
	List []ContentStatItem `json:"list"`
}

// 金币流水构成: 近 N 天按场景聚合收入/支出
type BalanceSceneItem struct {
	Scene   string  `json:"scene"`
	Income  float64 `json:"income"`  // direction=1 合计
	Expense float64 `json:"expense"` // direction=2 合计
}
type BalanceScenesReq struct {
	g.Meta `path:"/stats/balance-scenes" method:"get" tags:"Backend/Stats" summary:"金币流水构成"`
	Days   int `json:"days" v:"min:0|max:90#days 不合法|days 最大90"`
}
type BalanceScenesRes struct {
	List []BalanceSceneItem `json:"list"`
}
