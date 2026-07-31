package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	statsdomain "github.com/JarvanDante/my_service/internal/modules/stats/domain"
)

type statsRepo struct{}

// NewStatsRepo 返回 stats 领域仓储实现(只读聚合)。
func NewStatsRepo() statsdomain.Repository { return &statsRepo{} }

func (r *statsRepo) Overview(ctx context.Context, today, yesterday int) (*statsdomain.Overview, error) {
	o := &statsdomain.Overview{}
	// 用户: 总数 / 今日新增 / 昨日新增 / 今日活跃(DAU)
	one, err := g.DB().GetOne(ctx, `
		SELECT count(*)                                            AS total_users,
		       count(*) FILTER (WHERE register_date = ?)           AS today_new,
		       count(*) FILTER (WHERE register_date = ?)           AS yesterday_new,
		       count(*) FILTER (WHERE last_login_at >= current_date) AS today_active
		  FROM users`, today, yesterday)
	if err != nil {
		return nil, err
	}
	o.TotalUsers = one["total_users"].Int()
	o.TodayNew = one["today_new"].Int()
	o.YesterdayNew = one["yesterday_new"].Int()
	o.TodayActive = one["today_active"].Int()
	// 充值(已支付)
	one, err = g.DB().GetOne(ctx, `
		SELECT coalesce(sum(amount), 0)                                          AS total_amount,
		       count(*)                                                          AS total_orders,
		       coalesce(sum(amount) FILTER (WHERE pay_at >= current_date), 0)    AS today_amount,
		       count(*) FILTER (WHERE pay_at >= current_date)                    AS today_orders,
		       coalesce(sum(amount) FILTER (WHERE pay_at >= current_date - 1
		                                      AND pay_at < current_date), 0)     AS yest_amount
		  FROM recharge_order WHERE status = 1`)
	if err != nil {
		return nil, err
	}
	o.TotalPaidAmount = one["total_amount"].Float64()
	o.TotalPaidOrders = one["total_orders"].Int()
	o.TodayPaidAmount = one["today_amount"].Float64()
	o.TodayPaidOrders = one["today_orders"].Int()
	o.YestPaidAmount = one["yest_amount"].Float64()
	return o, nil
}

func (r *statsRepo) NewUsersByDay(ctx context.Context, startDate int) ([]statsdomain.DayCount, error) {
	all, err := g.DB().GetAll(ctx, `
		SELECT register_date AS d, count(*) AS cnt
		  FROM users WHERE register_date >= ?
		 GROUP BY register_date ORDER BY register_date`, startDate)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.DayCount, 0, len(all))
	for _, rec := range all {
		out = append(out, statsdomain.DayCount{Date: rec["d"].Int(), Count: rec["cnt"].Int()})
	}
	return out, nil
}

func (r *statsRepo) PaidOrdersByDay(ctx context.Context, startDay string) ([]statsdomain.DayAmount, error) {
	all, err := g.DB().GetAll(ctx, `
		SELECT to_char(pay_at, 'YYYYMMDD')::int AS d,
		       count(*)                          AS cnt,
		       coalesce(sum(amount), 0)          AS amt
		  FROM recharge_order
		 WHERE status = 1 AND pay_at >= ?::date
		 GROUP BY 1 ORDER BY 1`, startDay)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.DayAmount, 0, len(all))
	for _, rec := range all {
		out = append(out, statsdomain.DayAmount{
			Date: rec["d"].Int(), Orders: rec["cnt"].Int(), Amount: rec["amt"].Float64(),
		})
	}
	return out, nil
}

func (r *statsRepo) ChannelStats(ctx context.Context) ([]statsdomain.ChannelStat, error) {
	all, err := g.DB().GetAll(ctx, `
		SELECT channel_name AS ch, count(*) AS cnt, coalesce(sum(money_count), 0) AS amt
		  FROM users GROUP BY channel_name ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.ChannelStat, 0, len(all))
	for _, rec := range all {
		out = append(out, statsdomain.ChannelStat{
			Channel: rec["ch"].String(), UserCount: rec["cnt"].Int(), TotalRecharge: rec["amt"].Float64(),
		})
	}
	return out, nil
}
