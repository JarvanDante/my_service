// Package service 统计模块对外接口(B8)。
package service

import "context"

type OverviewDTO struct {
	TotalUsers      int
	TodayNew        int
	YesterdayNew    int
	TodayActive     int
	TodayPaidAmount float64
	TodayPaidOrders int
	YestPaidAmount  float64
	TotalPaidAmount float64
	TotalPaidOrders int
}

type UserTrendItem struct {
	Date     int // YYYYMMDD
	NewUsers int
}

type RechargeTrendItem struct {
	Date   int
	Orders int
	Amount float64
}

type ChannelStatDTO struct {
	Channel       string
	UserCount     int
	TotalRecharge float64
}

type IStats interface {
	Overview(ctx context.Context) (*OverviewDTO, error)
	UserTrend(ctx context.Context, days int) ([]UserTrendItem, error)
	RechargeTrend(ctx context.Context, days int) ([]RechargeTrendItem, error)
	Channels(ctx context.Context) ([]ChannelStatDTO, error)
	// ---- 扩展维度(分析页图表化) ----
	HourDist(ctx context.Context, days int) ([]HourDistItem, error)
	DeviceStats(ctx context.Context) ([]DeviceStatDTO, error)
	ContentStats(ctx context.Context) ([]ContentStatDTO, error)
	BalanceScenes(ctx context.Context, days int) ([]BalanceSceneDTO, error)
}

// ---------------- 扩展维度 ----------------

type HourDistItem struct {
	Hour      int
	Registers int
	Orders    int
}

type DeviceStatDTO struct {
	DeviceType string
	Count      int
}

type ContentStatDTO struct {
	MediaType int
	TypeName  string
	Online    int
	Pending   int
	Offline   int
	Views     int64
	Buys      int
	BuyAmount float64
}

type BalanceSceneDTO struct {
	Scene   string
	Income  float64
	Expense float64
}
