// Package service 商品兑换对外接口。
package service

import "context"

type GoodsDTO struct {
	Id        int64
	Name      string
	Cover     string
	Intro     string
	CostGold  float64
	Stock     int
	Exchanged int
	Rank      int
	Status    int
	CreatedAt string
}

type OrderDTO struct {
	Id        int64
	UserId    int64
	GoodsId   int64
	GoodsName string
	CostGold  float64
	CreatedAt string
}

type SaveInput struct {
	Id       int64 // 0=新增
	Name     string
	Cover    string
	Intro    string
	CostGold float64
	Stock    int
	Rank     int
	Status   int
}

type ListFilter struct {
	Status  int // -1=全部
	Keyword string
	UserId  int64
	GoodsId int64
	Page    int
	Size    int
}

type IRedeemGoods interface {
	FrontList(ctx context.Context, page, size int) ([]*GoodsDTO, int, error)
	// Exchange 事务: 余额条件扣款 + 库存条件递减 + 流水 + 记录, 全程防超卖/防透支。
	Exchange(ctx context.Context, userId, goodsId int64) (int64, error)
	History(ctx context.Context, userId int64, page, size int) ([]*OrderDTO, error)
	List(ctx context.Context, f ListFilter) ([]*GoodsDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	Orders(ctx context.Context, f ListFilter) ([]*OrderDTO, int, error)
}
