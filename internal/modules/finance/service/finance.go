// Package service 财务模块对外接口(B2)。
package service

import "context"

// ---------- 套餐 ----------

type RechargePackageDTO struct {
	Id     int64
	Name   string
	Amount float64
	Coin   float64
	Bonus  float64
	Sort   int
	Status int
}

type RechargePackageInput struct {
	Id     int64 // 更新时用
	Name   string
	Amount float64
	Coin   float64
	Bonus  float64
	Sort   int
	Status int
}

type VipPackageDTO struct {
	Id      int64
	Name    string
	Days    int
	Price   float64
	GroupId int64
	Sort    int
	Status  int
}

type VipPackageInput struct {
	Id      int64
	Name    string
	Days    int
	Price   float64
	GroupId int64
	Sort    int
	Status  int
}

// ---------- 订单 / 流水 ----------

type OrderListInput struct {
	OrderNo   string
	UserId    int64
	Status    int // 0全部 1待支付 2已支付 3已取消
	StartDate string
	EndDate   string
	Page      int
	Size      int
}

type OrderDTO struct {
	Id        int64
	OrderNo   string
	UserId    int64
	PackageId int64
	Amount    float64
	Coin      float64
	Status    int // 数据库原值: 0待支付 1已支付 -1取消
	PayAt     string
	CreatedAt string
}

type BalanceLogListInput struct {
	UserId    int64
	Scene     string
	Direction int
	StartDate string
	EndDate   string
	Page      int
	Size      int
}

type BalanceLogDTO struct {
	Id            int64
	UserId        int64
	Direction     int
	Scene         string
	Amount        float64
	BalanceBefore float64
	BalanceAfter  float64
	RefId         string
	Remark        string
	CreatedAt     string
}

type PageDTO[T any] struct {
	List  []*T
	Total int
	Page  int
	Size  int
}

type IFinance interface {
	// 充值套餐
	RechargePackages(ctx context.Context) ([]*RechargePackageDTO, error)
	CreateRechargePackage(ctx context.Context, in RechargePackageInput) (int64, error)
	UpdateRechargePackage(ctx context.Context, in RechargePackageInput) error
	DeleteRechargePackage(ctx context.Context, id int64) error
	// VIP 套餐
	VipPackages(ctx context.Context) ([]*VipPackageDTO, error)
	CreateVipPackage(ctx context.Context, in VipPackageInput) (int64, error)
	UpdateVipPackage(ctx context.Context, in VipPackageInput) error
	DeleteVipPackage(ctx context.Context, id int64) error
	// 订单 / 流水
	Orders(ctx context.Context, in OrderListInput) (*PageDTO[OrderDTO], error)
	BalanceLogs(ctx context.Context, in BalanceLogListInput) (*PageDTO[BalanceLogDTO], error)
	// 支付回调
	PayCallback(ctx context.Context, orderNo, sign string) error
}
