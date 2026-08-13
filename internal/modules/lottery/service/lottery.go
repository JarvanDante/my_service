// Package service 抽奖对外接口。
package service

import "context"

// 玩法类型(lottery_activity.lottery_type)
const (
	TypeVipDay  = 1 // 会员日抽奖
	TypeWelfare = 2 // 福利抽奖
)

// 消耗方式(lottery_activity.pay_type / lottery_history.pay_type)
const (
	PayFree = 1 // 免费次数
	PayGold = 2 // 金币
)

// 奖品类型(lottery_prize.type / lottery_history.prize_type)
const (
	PrizeGold   = 1 // 金币
	PrizeVip    = 2 // VIP天数
	PrizeCoupon = 3 // 优惠券
	PrizeGoods  = 4 // 实物
	PrizeThanks = 5 // 谢谢参与
)

// 中奖记录状态(lottery_history.status)
const (
	HistoryDone     = 1 // 已发放(金币/VIP/券/谢谢参与)
	HistoryWaitShip = 2 // 待发货(实物)
	HistoryShipped  = 3 // 已发货
)

// 收货单状态(prize_addr.delivery_status)
const (
	AddrWaitFill = 0 // 待用户填写
	AddrWaitShip = 1 // 待发货
	AddrShipped  = 2 // 已发货
)

func PrizeTypeText(t int) string {
	switch t {
	case PrizeGold:
		return "金币"
	case PrizeVip:
		return "VIP天数"
	case PrizeCoupon:
		return "优惠券"
	case PrizeGoods:
		return "实物"
	case PrizeThanks:
		return "谢谢参与"
	}
	return "未知"
}

type ActivityDTO struct {
	Id          int64
	Name        string
	LotteryType int
	PayType     int
	CostGold    float64
	DailyFree   int
	DailyLimit  int
	Notice      string
	Status      int
	CreatedAt   string
}

// PrizeDTO 奖品。Odds 只在后台 DTO 里回填, 前台控制器不允许下发(概率是商业机密)。
type PrizeDTO struct {
	Id          int64
	ActivityId  int64
	Name        string
	Cover       string
	Desc        string
	Type        int
	Amount      float64
	CouponTplId int64
	Odds        int
	Stock       int
	Awarded     int
	Rank        int
	Status      int
	CreatedAt   string
}

type HistoryDTO struct {
	Id          int64
	UserId      int64
	Nickname    string
	ActivityId  int64
	LotteryType int
	PayType     int
	CostGold    float64
	PrizeId     int64
	PrizeName   string
	PrizeType   int
	PrizeAmount float64
	Status      int
	Remark      string
	CreatedAt   string
}

// InfoDTO 前台活动首页数据(活动 + 奖品 + 我的次数/余额 + 全站跑马灯)。
type InfoDTO struct {
	Activity *ActivityDTO
	Prizes   []*PrizeDTO
	FreeLeft int     // 今日剩余免费次数
	DrawLeft int     // 今日剩余可抽次数(daily_limit=0 时为 -1 表示不限)
	Balance  float64 // 我的金币
	Marquee  []*HistoryDTO
	LoggedIn bool
}

// DrawDTO 一次抽奖的结果。
type DrawDTO struct {
	HistoryId   int64
	PrizeId     int64
	PrizeName   string
	PrizeType   int
	PrizeAmount float64
	PrizeCover  string
	PrizeDesc   string
	Status      int
	NeedAddr    bool // 实物奖: 前端要引导填收货地址
	Remark      string
	FreeLeft    int
	DrawLeft    int
	Balance     float64
	CostGold    float64
}

type AddrDTO struct {
	Id             int64
	HistoryId      int64
	UserId         int64
	Nickname       string
	PrizeName      string
	Receiver       string
	Phone          string
	Address        string
	DeliveryStatus int
	ExpressNo      string
	CreatedAt      string
}

type ActivityInput struct {
	Id          int64
	Name        string
	LotteryType int
	PayType     int
	CostGold    float64
	DailyFree   int
	DailyLimit  int
	Notice      string
	Status      int
}

type PrizeInput struct {
	Id          int64
	ActivityId  int64
	Name        string
	Cover       string
	Desc        string
	Type        int
	Amount      float64
	CouponTplId int64
	Odds        int
	Stock       int
	Rank        int
	Status      int
}

// HistoryFilter 后台中奖记录筛选, -1 表示不限。
type HistoryFilter struct {
	UserId      int64
	LotteryType int
	PrizeType   int
	Status      int
	Page        int
	Size        int
}

type AddrFilter struct {
	DeliveryStatus int // -1=全部
	UserId         int64
	Page           int
	Size           int
}

type ILottery interface {
	// Info 前台活动信息: 活动 + 奖品(不含 odds) + 我的次数与余额 + 最近中奖跑马灯。
	Info(ctx context.Context, userId int64, lotteryType int) (*InfoDTO, error)
	// Draw 抽奖(全程单事务: 校验 → 扣费 → 加权随机 → 库存条件递减 → 发奖 → 写记录)。
	Draw(ctx context.Context, userId int64, lotteryType, payType int) (*DrawDTO, error)
	// My 我的中奖记录。
	My(ctx context.Context, userId int64, lotteryType, page, size int) ([]*HistoryDTO, int, error)
	// FillAddr 填写实物奖收货地址(只能填自己的实物奖且 delivery_status=0)。
	FillAddr(ctx context.Context, userId, historyId int64, receiver, phone, address string) error
	// MyAddrs 我的收货单。
	MyAddrs(ctx context.Context, userId int64, page, size int) ([]*AddrDTO, int, error)

	Activities(ctx context.Context, status int) ([]*ActivityDTO, error)
	ActivityCreate(ctx context.Context, in ActivityInput) (int64, error)
	ActivityUpdate(ctx context.Context, in ActivityInput) error
	ActivityDelete(ctx context.Context, id int64) error

	Prizes(ctx context.Context, activityId int64) ([]*PrizeDTO, error)
	PrizeCreate(ctx context.Context, in PrizeInput) (int64, error)
	PrizeUpdate(ctx context.Context, in PrizeInput) error
	PrizeDelete(ctx context.Context, id int64) error

	Histories(ctx context.Context, f HistoryFilter) ([]*HistoryDTO, int, error)
	Addrs(ctx context.Context, f AddrFilter) ([]*AddrDTO, int, error)
	// Ship 后台发货: 条件更新 delivery_status 1→2 并把 history 置为已发货。
	Ship(ctx context.Context, addrId int64, expressNo string) error
}
