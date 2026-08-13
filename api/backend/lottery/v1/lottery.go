// Package v1 后台抽奖契约(活动/奖品 CRUD + 中奖记录 + 收货单发货)。
// 约定: 所有筛选参数一律 string 接收, 空串=全部(int 的零值分不清"没传"和"传了0")。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------------------------------------------------------------- 活动

type ActivityItem struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	LotteryType int     `json:"lottery_type"` // 1会员日 2福利
	PayType     int     `json:"pay_type"`     // 1仅免费次数 2金币
	CostGold    float64 `json:"cost_gold"`
	DailyFree   int     `json:"daily_free"`
	DailyLimit  int     `json:"daily_limit"`
	Notice      string  `json:"notice"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type ActivityListReq struct {
	g.Meta `path:"/lottery/activities" method:"get" tags:"Backend/Lottery" summary:"抽奖活动列表"`
	Status string `json:"status"`
}
type ActivityListRes struct {
	List []ActivityItem `json:"list"`
}

type ActivityCreateReq struct {
	g.Meta      `path:"/lottery/activities" method:"post" tags:"Backend/Lottery" summary:"新增抽奖活动"`
	Name        string  `json:"name"         v:"required#活动名必填"`
	LotteryType int     `json:"lottery_type" v:"required|in:1,2#玩法类型非法"`
	PayType     int     `json:"pay_type"     v:"in:1,2#消耗方式非法"`
	CostGold    float64 `json:"cost_gold"`
	DailyFree   int     `json:"daily_free"`
	DailyLimit  int     `json:"daily_limit"`
	Notice      string  `json:"notice"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type ActivityCreateRes struct {
	Id int64 `json:"id"`
}

type ActivityUpdateReq struct {
	g.Meta     `path:"/lottery/activities/{id}" method:"put" tags:"Backend/Lottery" summary:"更新抽奖活动"`
	Id         int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name       string  `json:"name"`
	PayType    int     `json:"pay_type"`
	CostGold   float64 `json:"cost_gold"`
	DailyFree  int     `json:"daily_free"`
	DailyLimit int     `json:"daily_limit"`
	Notice     string  `json:"notice"`
	Status     int     `json:"status" v:"in:0,1#状态非法"`
}
type ActivityUpdateRes struct{}

type ActivityDeleteReq struct {
	g.Meta `path:"/lottery/activities/{id}" method:"delete" tags:"Backend/Lottery" summary:"删除抽奖活动(连带奖品)"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type ActivityDeleteRes struct{}

// ---------------------------------------------------------------- 奖品

type PrizeItem struct {
	Id          int64   `json:"id"`
	ActivityId  int64   `json:"activity_id"`
	Name        string  `json:"name"`
	Cover       string  `json:"cover"`
	Desc        string  `json:"desc"`
	Type        int     `json:"type"` // 1金币 2VIP天数 3优惠券 4实物 5谢谢参与
	Amount      float64 `json:"amount"`
	CouponTplId int64   `json:"coupon_tpl_id"`
	Odds        int     `json:"odds"` // 整数权重, 只在后台可见
	Stock       int     `json:"stock"`
	Awarded     int     `json:"awarded"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type PrizeListReq struct {
	g.Meta     `path:"/lottery/prizes" method:"get" tags:"Backend/Lottery" summary:"抽奖奖品列表"`
	ActivityId string `json:"activity_id"`
}
type PrizeListRes struct {
	List []PrizeItem `json:"list"`
}

type PrizeCreateReq struct {
	g.Meta      `path:"/lottery/prizes" method:"post" tags:"Backend/Lottery" summary:"新增抽奖奖品"`
	ActivityId  int64   `json:"activity_id" v:"required|min:1#活动ID必填"`
	Name        string  `json:"name"        v:"required#奖品名必填"`
	Cover       string  `json:"cover"`
	Desc        string  `json:"desc"`
	Type        int     `json:"type" v:"required|in:1,2,3,4,5#奖品类型非法"`
	Amount      float64 `json:"amount"`
	CouponTplId int64   `json:"coupon_tpl_id"`
	Odds        int     `json:"odds"`
	Stock       int     `json:"stock"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type PrizeCreateRes struct {
	Id int64 `json:"id"`
}

type PrizeUpdateReq struct {
	g.Meta      `path:"/lottery/prizes/{id}" method:"put" tags:"Backend/Lottery" summary:"更新抽奖奖品"`
	Id          int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name        string  `json:"name"`
	Cover       string  `json:"cover"`
	Desc        string  `json:"desc"`
	Type        int     `json:"type"`
	Amount      float64 `json:"amount"`
	CouponTplId int64   `json:"coupon_tpl_id"`
	Odds        int     `json:"odds"`
	Stock       int     `json:"stock"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type PrizeUpdateRes struct{}

type PrizeDeleteReq struct {
	g.Meta `path:"/lottery/prizes/{id}" method:"delete" tags:"Backend/Lottery" summary:"删除抽奖奖品"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type PrizeDeleteRes struct{}

// ---------------------------------------------------------------- 中奖记录

type HistoryItem struct {
	Id          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
	Nickname    string  `json:"nickname"`
	ActivityId  int64   `json:"activity_id"`
	LotteryType int     `json:"lottery_type"`
	PayType     int     `json:"pay_type"`
	CostGold    float64 `json:"cost_gold"`
	PrizeId     int64   `json:"prize_id"`
	PrizeName   string  `json:"prize_name"`
	PrizeType   int     `json:"prize_type"`
	PrizeText   string  `json:"prize_text"`
	PrizeAmount float64 `json:"prize_amount"`
	Status      int     `json:"status"`
	Remark      string  `json:"remark"`
	CreatedAt   string  `json:"created_at"`
}

// HistoryListReq 中奖记录(筛选项全部 string, 空=全部)。
type HistoryListReq struct {
	g.Meta      `path:"/lottery/histories" method:"get" tags:"Backend/Lottery" summary:"中奖记录列表"`
	UserId      string `json:"user_id"`
	LotteryType string `json:"lottery_type"`
	PrizeType   string `json:"prize_type"`
	Status      string `json:"status"`
	Page        int    `json:"page"`
	Size        int    `json:"size"`
}
type HistoryListRes struct {
	List  []HistoryItem `json:"list"`
	Total int           `json:"total"`
}

// ---------------------------------------------------------------- 收货单

type AddrItem struct {
	Id             int64  `json:"id"`
	HistoryId      int64  `json:"history_id"`
	UserId         int64  `json:"user_id"`
	Nickname       string `json:"nickname"`
	PrizeName      string `json:"prize_name"`
	Receiver       string `json:"receiver"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	DeliveryStatus int    `json:"delivery_status"` // 0待填写 1待发货 2已发货
	ExpressNo      string `json:"express_no"`
	CreatedAt      string `json:"created_at"`
}

type AddrListReq struct {
	g.Meta         `path:"/lottery/addresses" method:"get" tags:"Backend/Lottery" summary:"实物奖收货单列表"`
	DeliveryStatus string `json:"delivery_status"`
	UserId         string `json:"user_id"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}
type AddrListRes struct {
	List  []AddrItem `json:"list"`
	Total int        `json:"total"`
}

type ShipReq struct {
	g.Meta    `path:"/lottery/addresses/{id}/ship" method:"post" tags:"Backend/Lottery" summary:"实物奖发货"`
	Id        int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	ExpressNo string `json:"express_no" v:"required#快递单号必填"`
}
type ShipRes struct{}
