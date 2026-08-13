// Code maintained manually (抽奖活动 / 奖品 / 中奖记录 / 实物收货单).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// LotteryActivity 抽奖活动。lottery_type: 1会员日 2福利; pay_type: 1仅免费次数 2金币。
type LotteryActivity struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	Name        string      `json:"name"        orm:"name"`
	LotteryType int         `json:"lotteryType" orm:"lottery_type"`
	PayType     int         `json:"payType"     orm:"pay_type"`
	CostGold    float64     `json:"costGold"    orm:"cost_gold"`
	DailyFree   int         `json:"dailyFree"   orm:"daily_free"`
	DailyLimit  int         `json:"dailyLimit"  orm:"daily_limit"`
	Notice      string      `json:"notice"      orm:"notice"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

// LotteryPrize 抽奖奖品。type: 1金币 2VIP天数 3优惠券 4实物 5谢谢参与。
// Odds 是整数权重(前缀和抽取), 不是百分比; Stock=-1 表示不限量。
type LotteryPrize struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	ActivityId  int64       `json:"activityId"  orm:"activity_id"`
	Name        string      `json:"name"        orm:"name"`
	Cover       string      `json:"cover"       orm:"cover"`
	Desc        string      `json:"desc"        orm:"desc"`
	Type        int         `json:"type"        orm:"type"`
	Amount      float64     `json:"amount"      orm:"amount"`
	CouponTplId int64       `json:"couponTplId" orm:"coupon_tpl_id"`
	Odds        int         `json:"odds"        orm:"odds"`
	Stock       int         `json:"stock"       orm:"stock"`
	Awarded     int         `json:"awarded"     orm:"awarded"`
	Rank        int         `json:"rank"        orm:"rank"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

// LotteryHistory 中奖记录(与扣费同事务写入)。status: 1已发放 2待发货 3已发货。
type LotteryHistory struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	UserId      int64       `json:"userId"      orm:"user_id"`
	ActivityId  int64       `json:"activityId"  orm:"activity_id"`
	LotteryType int         `json:"lotteryType" orm:"lottery_type"`
	PayType     int         `json:"payType"     orm:"pay_type"`
	CostGold    float64     `json:"costGold"    orm:"cost_gold"`
	PrizeId     int64       `json:"prizeId"     orm:"prize_id"`
	PrizeName   string      `json:"prizeName"   orm:"prize_name"`
	PrizeType   int         `json:"prizeType"   orm:"prize_type"`
	PrizeAmount float64     `json:"prizeAmount" orm:"prize_amount"`
	Status      int         `json:"status"      orm:"status"`
	Remark      string      `json:"remark"      orm:"remark"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

// PrizeAddr 实物奖收货单。delivery_status: 0待填写 1待发货 2已发货。
type PrizeAddr struct {
	Id             int64       `json:"id"             orm:"id"`
	SiteId         int64       `json:"siteId"         orm:"site_id"`
	HistoryId      int64       `json:"historyId"      orm:"history_id"`
	UserId         int64       `json:"userId"         orm:"user_id"`
	Receiver       string      `json:"receiver"       orm:"receiver"`
	Phone          string      `json:"phone"          orm:"phone"`
	Address        string      `json:"address"        orm:"address"`
	DeliveryStatus int         `json:"deliveryStatus" orm:"delivery_status"`
	ExpressNo      string      `json:"expressNo"      orm:"express_no"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"`
}
