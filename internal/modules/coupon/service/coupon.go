// Package service 优惠券对外接口。
package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// 券类型 / 场景 / 状态
const (
	TypeCash     = 1 // 抵用券(直接减面额)
	TypeDiscount = 2 // 折扣券(按折扣率算, 可设最高抵扣)

	SceneRecharge = 1 // 充值
	SceneContent  = 2 // 内容购买
	SceneAll      = 3 // 通用

	StatusUnused = 1
	StatusUsed   = 2
	StatusExpire = 3
)

func StatusText(s int) string {
	switch s {
	case StatusUnused:
		return "未使用"
	case StatusUsed:
		return "已使用"
	case StatusExpire:
		return "已过期"
	}
	return "未知"
}

type TplDTO struct {
	Id        int64
	Name      string
	Type      int
	Scene     int
	FaceValue float64
	Discount  int
	Threshold float64
	MaxDeduct float64
	Total     int
	Issued    int
	PerLimit  int
	ExpireDay int
	Status    int
	CreatedAt string
	Received  bool // 前台用: 当前用户是否已领满
}

type UserCouponDTO struct {
	Id        int64
	UserId    int64
	TplId     int64
	Name      string
	Type      int
	Scene     int
	FaceValue float64
	Discount  int
	Threshold float64
	MaxDeduct float64
	Status    int
	RefId     string
	ExpireAt  string
	UsedAt    string
	CreatedAt string
	Deduct    float64 // 针对某订单金额算出的抵扣额(仅 Available 接口填充)
}

type TplInput struct {
	Id        int64
	Name      string
	Type      int
	Scene     int
	FaceValue float64
	Discount  int
	Threshold float64
	MaxDeduct float64
	Total     int
	PerLimit  int
	ExpireDay int
	Status    int
}

type ListFilter struct {
	Status  int // -1=全部
	Keyword string
	UserId  int64
	TplId   int64
	Page    int
	Size    int
}

type ICoupon interface {
	// Tpls 前台可领券列表(带"是否已领满"标记)。
	Tpls(ctx context.Context, userId int64) ([]*TplDTO, error)
	// Receive 领券: 总量条件递增(防超发) + 每人限领校验。
	Receive(ctx context.Context, userId, tplId int64) (int64, error)
	// My 我的券(顺带把已过期的未使用券刷成 3)。
	My(ctx context.Context, userId int64, status, page, size int) ([]*UserCouponDTO, int, error)
	// Available 给定订单金额/场景, 返回可用券(含抵扣额)与最优券。
	Available(ctx context.Context, userId int64, scene int, amount float64) ([]*UserCouponDTO, int64, float64, error)
	// UseInTx 核销(供下单链路在自己的事务里调用): 条件更新 1→2, 返回实际抵扣额。
	UseInTx(ctx context.Context, tx gdb.TX, userId, couponId int64, scene int, amount float64, refId string) (float64, error)

	List(ctx context.Context, f ListFilter) ([]*TplDTO, int, error)
	Create(ctx context.Context, in TplInput) (int64, error)
	Update(ctx context.Context, in TplInput) error
	Delete(ctx context.Context, id int64) error
	Grant(ctx context.Context, tplId int64, userIds []int64) (int, int, []string)
	Users(ctx context.Context, f ListFilter) ([]*UserCouponDTO, int, error)
}
