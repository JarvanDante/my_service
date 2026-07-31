// Package logic 统计模块业务实现(B8)。
// 说明: 活跃/留存的历史趋势需要登录流水表, 当前仅有 users.last_login_at(只存最近一次),
// TODO: 加 user_login_log 表后补 DAU 趋势与留存分析。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/modules/stats/domain"
	"github.com/JarvanDante/my_service/internal/modules/stats/service"
)

type sStats struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IStats { return &sStats{repo: repo} }

func (s *sStats) Overview(ctx context.Context) (*service.OverviewDTO, error) {
	now := gtime.Now()
	today := gconv.Int(now.Format("Ymd"))
	yesterday := gconv.Int(now.AddDate(0, 0, -1).Format("Ymd"))
	o, err := s.repo.Overview(ctx, today, yesterday)
	if err != nil {
		return nil, err
	}
	return &service.OverviewDTO{
		TotalUsers: o.TotalUsers, TodayNew: o.TodayNew, YesterdayNew: o.YesterdayNew,
		TodayActive:     o.TodayActive,
		TodayPaidAmount: o.TodayPaidAmount, TodayPaidOrders: o.TodayPaidOrders,
		YestPaidAmount:  o.YestPaidAmount,
		TotalPaidAmount: o.TotalPaidAmount, TotalPaidOrders: o.TotalPaidOrders,
	}, nil
}

// lastNDays 生成含今天在内的最近 N 天 YYYYMMDD 序列(升序)。
func lastNDays(n int) []int {
	now := gtime.Now()
	out := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		out = append(out, gconv.Int(now.AddDate(0, 0, -i).Format("Ymd")))
	}
	return out
}

func checkDays(days int) (int, error) {
	if days == 0 {
		return 7, nil
	}
	if days < 1 || days > 90 {
		return 0, gerror.New("days 须在 1~90")
	}
	return days, nil
}

func (s *sStats) UserTrend(ctx context.Context, days int) ([]service.UserTrendItem, error) {
	days, err := checkDays(days)
	if err != nil {
		return nil, err
	}
	keys := lastNDays(days)
	rows, err := s.repo.NewUsersByDay(ctx, keys[0])
	if err != nil {
		return nil, err
	}
	byDate := make(map[int]int, len(rows))
	for _, r := range rows {
		byDate[r.Date] = r.Count
	}
	out := make([]service.UserTrendItem, 0, days)
	for _, d := range keys {
		out = append(out, service.UserTrendItem{Date: d, NewUsers: byDate[d]})
	}
	return out, nil
}

func (s *sStats) RechargeTrend(ctx context.Context, days int) ([]service.RechargeTrendItem, error) {
	days, err := checkDays(days)
	if err != nil {
		return nil, err
	}
	keys := lastNDays(days)
	startDay := gtime.Now().AddDate(0, 0, -(days - 1)).Format("Y-m-d")
	rows, err := s.repo.PaidOrdersByDay(ctx, startDay)
	if err != nil {
		return nil, err
	}
	type oa struct {
		orders int
		amount float64
	}
	byDate := make(map[int]oa, len(rows))
	for _, r := range rows {
		byDate[r.Date] = oa{orders: r.Orders, amount: r.Amount}
	}
	out := make([]service.RechargeTrendItem, 0, days)
	for _, d := range keys {
		v := byDate[d]
		out = append(out, service.RechargeTrendItem{Date: d, Orders: v.orders, Amount: v.amount})
	}
	return out, nil
}

func (s *sStats) Channels(ctx context.Context) ([]service.ChannelStatDTO, error) {
	rows, err := s.repo.ChannelStats(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]service.ChannelStatDTO, 0, len(rows))
	for _, r := range rows {
		ch := r.Channel
		if ch == "" {
			ch = "(未知)"
		}
		out = append(out, service.ChannelStatDTO{
			Channel: ch, UserCount: r.UserCount, TotalRecharge: r.TotalRecharge,
		})
	}
	return out, nil
}
