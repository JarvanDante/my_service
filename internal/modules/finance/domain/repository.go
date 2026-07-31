// Package domain 财务模块领域层(B2: 套餐/订单/流水/回调)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

// OrderFilter 充值订单筛选。
type OrderFilter struct {
	OrderNo   string
	UserId    int64
	Status    int    // 0全部 1待支付 2已支付 3已取消
	StartDate string // created_at >= (YYYY-MM-DD)
	EndDate   string // created_at <  次日
}

// BalanceLogFilter 全站余额流水筛选。
type BalanceLogFilter struct {
	UserId    int64
	Scene     string
	Direction int // 0全部 1收入 2支出
	StartDate string
	EndDate   string
}

type Repository interface {
	// 充值套餐 CRUD(后台可见含下架)
	ListRechargePackages(ctx context.Context) ([]*entity.RechargePackage, error)
	CreateRechargePackage(ctx context.Context, name string, amount, coin, bonus float64, sort, status int) (int64, error)
	UpdateRechargePackage(ctx context.Context, id int64, name string, amount, coin, bonus float64, sort, status int) error
	DeleteRechargePackage(ctx context.Context, id int64) error

	// VIP 套餐 CRUD
	ListVipPackages(ctx context.Context) ([]*entity.VipPackage, error)
	CreateVipPackage(ctx context.Context, name string, days int, price float64, groupId int64, sort, status int) (int64, error)
	UpdateVipPackage(ctx context.Context, id int64, name string, days int, price float64, groupId int64, sort, status int) error
	DeleteVipPackage(ctx context.Context, id int64) error

	// 订单 / 流水
	ListOrders(ctx context.Context, f OrderFilter, page, size int) ([]*entity.RechargeOrder, int, error)
	ListBalanceLogs(ctx context.Context, f BalanceLogFilter, page, size int) ([]*entity.UserBalanceLog, int, error)

	// 支付回调: 置已支付 + 用户到账 + 写流水(幂等: 已处理返回 ErrOrderHandled)
	MarkOrderPaid(ctx context.Context, orderNo string) error
}
