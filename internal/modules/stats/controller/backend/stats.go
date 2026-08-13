// Package backend 后台统计控制器(B8)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/stats/v1"
	"github.com/JarvanDante/my_service/internal/modules/stats/service"
)

type Controller struct{ stats service.IStats }

func New(svc service.IStats) *Controller { return &Controller{stats: svc} }

// Overview 数据概览。
func (c *Controller) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	o, err := c.stats.Overview(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.OverviewRes{
		TotalUsers: o.TotalUsers, TodayNew: o.TodayNew, YesterdayNew: o.YesterdayNew,
		TodayActive:     o.TodayActive,
		TodayPaidAmount: o.TodayPaidAmount, TodayPaidOrders: o.TodayPaidOrders,
		YestPaidAmount:  o.YestPaidAmount,
		TotalPaidAmount: o.TotalPaidAmount, TotalPaidOrders: o.TotalPaidOrders,
	}, nil
}

// UserTrend 注册趋势。
func (c *Controller) UserTrend(ctx context.Context, req *v1.UserTrendReq) (res *v1.UserTrendRes, err error) {
	list, err := c.stats.UserTrend(ctx, req.Days)
	if err != nil {
		return nil, err
	}
	items := make([]v1.UserTrendItem, 0, len(list))
	for _, it := range list {
		items = append(items, v1.UserTrendItem{Date: it.Date, NewUsers: it.NewUsers})
	}
	return &v1.UserTrendRes{List: items}, nil
}

// RechargeTrend 充值趋势。
func (c *Controller) RechargeTrend(ctx context.Context, req *v1.RechargeTrendReq) (res *v1.RechargeTrendRes, err error) {
	list, err := c.stats.RechargeTrend(ctx, req.Days)
	if err != nil {
		return nil, err
	}
	items := make([]v1.RechargeTrendItem, 0, len(list))
	for _, it := range list {
		items = append(items, v1.RechargeTrendItem{Date: it.Date, Orders: it.Orders, Amount: it.Amount})
	}
	return &v1.RechargeTrendRes{List: items}, nil
}

// ChannelStats 渠道分析。
func (c *Controller) ChannelStats(ctx context.Context, req *v1.ChannelStatsReq) (res *v1.ChannelStatsRes, err error) {
	list, err := c.stats.Channels(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.ChannelStatItem, 0, len(list))
	for _, it := range list {
		items = append(items, v1.ChannelStatItem{
			Channel: it.Channel, UserCount: it.UserCount, TotalRecharge: it.TotalRecharge,
		})
	}
	return &v1.ChannelStatsRes{List: items}, nil
}

// ---------------- 扩展维度(分析页图表化) ----------------

func (c *Controller) HourDist(ctx context.Context, req *v1.HourDistReq) (res *v1.HourDistRes, err error) {
	list, err := c.stats.HourDist(ctx, req.Days)
	if err != nil {
		return nil, err
	}
	res = &v1.HourDistRes{List: make([]v1.HourDistItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.HourDistItem{Hour: d.Hour, Registers: d.Registers, Orders: d.Orders})
	}
	return res, nil
}

func (c *Controller) DeviceStats(ctx context.Context, req *v1.DeviceStatsReq) (res *v1.DeviceStatsRes, err error) {
	list, err := c.stats.DeviceStats(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.DeviceStatsRes{List: make([]v1.DeviceStatItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.DeviceStatItem{DeviceType: d.DeviceType, Count: d.Count})
	}
	return res, nil
}

func (c *Controller) ContentStats(ctx context.Context, req *v1.ContentStatsReq) (res *v1.ContentStatsRes, err error) {
	list, err := c.stats.ContentStats(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.ContentStatsRes{List: make([]v1.ContentStatItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.ContentStatItem{
			MediaType: d.MediaType, TypeName: d.TypeName,
			Online: d.Online, Pending: d.Pending, Offline: d.Offline,
			Views: d.Views, Buys: d.Buys, BuyAmount: d.BuyAmount,
		})
	}
	return res, nil
}

func (c *Controller) BalanceScenes(ctx context.Context, req *v1.BalanceScenesReq) (res *v1.BalanceScenesRes, err error) {
	list, err := c.stats.BalanceScenes(ctx, req.Days)
	if err != nil {
		return nil, err
	}
	res = &v1.BalanceScenesRes{List: make([]v1.BalanceSceneItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.BalanceSceneItem{Scene: d.Scene, Income: d.Income, Expense: d.Expense})
	}
	return res, nil
}
