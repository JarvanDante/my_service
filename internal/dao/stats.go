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

// ---------------- 扩展维度(分析页图表化) ----------------

// HourDist 近 N 天注册/支付按小时聚合。两类事件各自 group by 后在内存合并,
// 比一条 FULL JOIN 直观, 数据量也就 24 行。
func (r *statsRepo) HourDist(ctx context.Context, startDay string) ([]statsdomain.HourCount, error) {
	regs, err := g.DB().GetAll(ctx, `
		SELECT extract(hour FROM register_at)::int AS h, count(*) AS cnt
		  FROM users WHERE register_at >= ?::timestamptz
		 GROUP BY 1`, startDay)
	if err != nil {
		return nil, err
	}
	orders, err := g.DB().GetAll(ctx, `
		SELECT extract(hour FROM pay_at)::int AS h, count(*) AS cnt
		  FROM recharge_order WHERE status = 1 AND pay_at >= ?::timestamptz
		 GROUP BY 1`, startDay)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.HourCount, 24)
	for i := range out {
		out[i].Hour = i
	}
	for _, row := range regs {
		if h := row["h"].Int(); h >= 0 && h < 24 {
			out[h].Registers = row["cnt"].Int()
		}
	}
	for _, row := range orders {
		if h := row["h"].Int(); h >= 0 && h < 24 {
			out[h].Orders = row["cnt"].Int()
		}
	}
	return out, nil
}

func (r *statsRepo) DeviceStats(ctx context.Context) ([]statsdomain.DeviceStat, error) {
	all, err := g.DB().GetAll(ctx, `
		SELECT coalesce(nullif(device_type, ''), '(未知)') AS dt, count(*) AS cnt
		  FROM users GROUP BY 1 ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.DeviceStat, 0, len(all))
	for _, row := range all {
		out = append(out, statsdomain.DeviceStat{DeviceType: row["dt"].String(), Count: row["cnt"].Int()})
	}
	return out, nil
}

// ContentStats 五类内容各查一次状态分布, 购买数据一次带出(按 media_type 分组)。
// 状态编码: video/comics/novel/photo 用 0待 1上 2下; post 用 0待审 1通过 2拒绝 3软删(3 归入下架口径)。
func (r *statsRepo) ContentStats(ctx context.Context) ([]statsdomain.ContentStat, error) {
	type src struct {
		mediaType int
		table     string
		viewCol   string // 没有观看数的表传空
	}
	srcs := []src{
		{1, "video", ""}, // video 表没有 view_count 列
		{2, "post", "view_count"},
		{3, "comics", "view_count"},
		{4, "novel", "view_count"},
		{5, "photo_album", "view_count"},
	}
	out := make([]statsdomain.ContentStat, 0, len(srcs))
	for _, s := range srcs {
		viewExpr := "0"
		if s.viewCol != "" {
			viewExpr = "coalesce(sum(" + s.viewCol + "), 0)"
		}
		one, err := g.DB().GetOne(ctx, `
			SELECT count(*) FILTER (WHERE status = 1) AS online,
			       count(*) FILTER (WHERE status = 0) AS pending,
			       count(*) FILTER (WHERE status >= 2) AS offline,
			       `+viewExpr+` AS views
			  FROM `+s.table)
		if err != nil {
			return nil, err
		}
		out = append(out, statsdomain.ContentStat{
			MediaType: s.mediaType,
			Online:    one["online"].Int(),
			Pending:   one["pending"].Int(),
			Offline:   one["offline"].Int(),
			Views:     one["views"].Int64(),
		})
	}
	buys, err := g.DB().GetAll(ctx, `
		SELECT media_type, count(*) AS cnt, coalesce(sum(amount), 0) AS amt
		  FROM content_purchase GROUP BY media_type`)
	if err != nil {
		return nil, err
	}
	for _, row := range buys {
		mt := row["media_type"].Int()
		for i := range out {
			if out[i].MediaType == mt {
				out[i].Buys = row["cnt"].Int()
				out[i].BuyAmount = row["amt"].Float64()
			}
		}
	}
	return out, nil
}

func (r *statsRepo) BalanceScenes(ctx context.Context, startDay string) ([]statsdomain.BalanceScene, error) {
	all, err := g.DB().GetAll(ctx, `
		SELECT scene,
		       coalesce(sum(amount) FILTER (WHERE direction = 1), 0) AS income,
		       coalesce(sum(amount) FILTER (WHERE direction = 2), 0) AS expense
		  FROM user_balance_log
		 WHERE created_at >= ?::timestamptz
		 GROUP BY scene ORDER BY greatest(
		       coalesce(sum(amount) FILTER (WHERE direction = 1), 0),
		       coalesce(sum(amount) FILTER (WHERE direction = 2), 0)) DESC`, startDay)
	if err != nil {
		return nil, err
	}
	out := make([]statsdomain.BalanceScene, 0, len(all))
	for _, row := range all {
		out = append(out, statsdomain.BalanceScene{
			Scene: row["scene"].String(), Income: row["income"].Float64(), Expense: row["expense"].Float64(),
		})
	}
	return out, nil
}
