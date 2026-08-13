// Package service 钱包对外接口。
package service

import "context"

type SummaryDTO struct {
	Balance   float64
	Frozen    float64
	TotalIn   float64
	TotalOut  float64
	Withdrawn float64
}

type WaterDTO struct {
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

type WaterFilter struct {
	UserId    int64 // 0=全部(后台)
	Direction int   // -1=全部
	Scene     string
	Page      int
	Size      int
}

type IWallet interface {
	// Summary 前台: 余额 + 在途冻结 + 累计收支。
	Summary(ctx context.Context, userId int64) (*SummaryDTO, error)
	// Waters 流水分页(前台传 userId, 后台可传 0 查全站)。
	Waters(ctx context.Context, f WaterFilter) ([]*WaterDTO, int, error)
	// Adjust 后台人工调账: amount>0 加币, amount<0 扣币(条件更新防扣成负数), 必写流水。
	Adjust(ctx context.Context, adminId, userId int64, amount float64, remark string) (float64, error)
}
