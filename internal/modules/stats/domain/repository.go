// Package domain 统计模块领域层(B8, 只读聚合)。
package domain

import "context"

// Overview 概览数据。
type Overview struct {
	TotalUsers      int
	TodayNew        int
	YesterdayNew    int
	TodayActive     int // DAU: 今日有登录
	TodayPaidAmount float64
	TodayPaidOrders int
	YestPaidAmount  float64
	TotalPaidAmount float64
	TotalPaidOrders int
}

// DayCount 按日计数(date=YYYYMMDD)。
type DayCount struct {
	Date  int
	Count int
}

// DayAmount 按日金额+单数(date=YYYYMMDD)。
type DayAmount struct {
	Date   int
	Orders int
	Amount float64
}

// ChannelStat 渠道分析。
type ChannelStat struct {
	Channel       string
	UserCount     int
	TotalRecharge float64 // 该渠道用户累计充值(users.money_count 汇总)
}

type Repository interface {
	Overview(ctx context.Context, today, yesterday int) (*Overview, error)
	NewUsersByDay(ctx context.Context, startDate int) ([]DayCount, error)      // register_date >= startDate 分组
	PaidOrdersByDay(ctx context.Context, startDay string) ([]DayAmount, error) // pay_at >= startDay 分组
	ChannelStats(ctx context.Context) ([]ChannelStat, error)
}
