package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	findomain "github.com/JarvanDante/my_service/internal/modules/finance/domain"
)

type financeRepo struct{}

// NewFinanceRepo 返回 finance 领域仓储实现。
func NewFinanceRepo() findomain.Repository { return &financeRepo{} }

// ---------- 充值套餐 ----------

func (r *financeRepo) ListRechargePackages(ctx context.Context) ([]*entity.RechargePackage, error) {
	var list []*entity.RechargePackage
	err := g.Model("recharge_package").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *financeRepo) CreateRechargePackage(ctx context.Context, name string, amount, coin, bonus float64, sort, status int) (int64, error) {
	res, err := g.Model("recharge_package").Ctx(ctx).Data(g.Map{
		"name": name, "amount": amount, "coin": coin, "bonus": bonus, "sort": sort, "status": status,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *financeRepo) UpdateRechargePackage(ctx context.Context, id int64, name string, amount, coin, bonus float64, sort, status int) error {
	_, err := g.Model("recharge_package").Ctx(ctx).Where("id", id).Data(g.Map{
		"name": name, "amount": amount, "coin": coin, "bonus": bonus, "sort": sort, "status": status,
	}).Update()
	return err
}

func (r *financeRepo) DeleteRechargePackage(ctx context.Context, id int64) error {
	_, err := g.Model("recharge_package").Ctx(ctx).Where("id", id).Delete()
	return err
}

// ---------- VIP 套餐 ----------

func (r *financeRepo) ListVipPackages(ctx context.Context) ([]*entity.VipPackage, error) {
	var list []*entity.VipPackage
	err := g.Model("vip_package").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *financeRepo) CreateVipPackage(ctx context.Context, name string, days int, price float64, groupId int64, sort, status int) (int64, error) {
	res, err := g.Model("vip_package").Ctx(ctx).Data(g.Map{
		"name": name, "days": days, "price": price, "group_id": groupId, "sort": sort, "status": status,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *financeRepo) UpdateVipPackage(ctx context.Context, id int64, name string, days int, price float64, groupId int64, sort, status int) error {
	_, err := g.Model("vip_package").Ctx(ctx).Where("id", id).Data(g.Map{
		"name": name, "days": days, "price": price, "group_id": groupId, "sort": sort, "status": status,
	}).Update()
	return err
}

func (r *financeRepo) DeleteVipPackage(ctx context.Context, id int64) error {
	_, err := g.Model("vip_package").Ctx(ctx).Where("id", id).Delete()
	return err
}

// ---------- 订单 / 流水 ----------

func (r *financeRepo) ListOrders(ctx context.Context, f findomain.OrderFilter, page, size int) ([]*entity.RechargeOrder, int, error) {
	m := g.Model("recharge_order").Ctx(ctx)
	if f.OrderNo != "" {
		m = m.Where("order_no", f.OrderNo)
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	switch f.Status { // 0全部 1待支付 2已支付 3已取消
	case 1:
		m = m.Where("status", 0)
	case 2:
		m = m.Where("status", 1)
	case 3:
		m = m.Where("status", -1)
	}
	if f.StartDate != "" {
		m = m.Where("created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		m = m.Where("created_at < ?::date + interval '1 day'", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RechargeOrder
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *financeRepo) ListBalanceLogs(ctx context.Context, f findomain.BalanceLogFilter, page, size int) ([]*entity.UserBalanceLog, int, error) {
	m := g.Model("user_balance_log").Ctx(ctx)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Scene != "" {
		m = m.Where("scene", f.Scene)
	}
	if f.Direction == 1 || f.Direction == 2 {
		m = m.Where("direction", f.Direction)
	}
	if f.StartDate != "" {
		m = m.Where("created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		m = m.Where("created_at < ?::date + interval '1 day'", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserBalanceLog
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// ---------- 支付回调 ----------

// MarkOrderPaid 事务: 锁订单 -> 校验待支付 -> 置已支付 -> 用户到账 -> 写流水。幂等。
func (r *financeRepo) MarkOrderPaid(ctx context.Context, orderNo string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var order *entity.RechargeOrder
		if err := tx.Model("recharge_order").Ctx(ctx).Where("order_no", orderNo).LockUpdate().Scan(&order); err != nil {
			return err
		}
		if order == nil {
			return gerror.New("订单不存在")
		}
		if order.Status == 1 {
			return gerror.New("订单已处理")
		}
		if order.Status != 0 {
			return gerror.New("订单已取消, 不能支付")
		}
		if _, err := tx.Model("recharge_order").Ctx(ctx).Where("id", order.Id).Data(g.Map{
			"status": 1, "pay_at": gtime.Now(),
		}).Update(); err != nil {
			return err
		}
		// 用户到账(锁用户行, 记录前后值)
		var u *entity.Users
		if err := tx.Model("users").Ctx(ctx).Where("id", order.UserId).LockUpdate().Scan(&u); err != nil {
			return err
		}
		if u == nil {
			return gerror.New("订单用户不存在")
		}
		before := u.Balance
		after := before + order.Coin
		if _, err := tx.Model("users").Ctx(ctx).Where("id", u.Id).Data(g.Map{
			"balance":     after,
			"money_count": &gdb.Counter{Field: "money_count", Value: order.Amount},
			"has_buy":     1,
		}).Update(); err != nil {
			return err
		}
		_, err := tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
			"user_id": u.Id, "direction": 1, "scene": "recharge",
			"amount": order.Coin, "balance_before": before, "balance_after": after,
			"ref_id": order.OrderNo, "remark": "充值到账",
		}).Insert()
		return err
	})
}
