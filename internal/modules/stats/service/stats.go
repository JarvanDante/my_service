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
}
