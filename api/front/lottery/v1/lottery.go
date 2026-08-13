// Package v1 前台抽奖契约(移植自 tianbi lotteryapi)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// PrizeItem 奖品(前台展示用)。
//
// 注意: **绝对不下发 odds(中奖权重)**。概率是商业机密, 一旦下发:
//   - 竞品可以直接抄走整套奖池配置;
//   - 用户能算出期望值, 运营调概率会被逐帧对比后挂到社交平台;
//   - 想作弊的人能精确知道哪个奖最值得刷。
//
// 所以 service.PrizeDTO 里虽然有 Odds, 前台控制器不会把它拷进这个结构。
type PrizeItem struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Cover  string  `json:"cover"`
	Desc   string  `json:"desc"`
	Type   int     `json:"type"` // 1金币 2VIP天数 3优惠券 4实物 5谢谢参与
	Amount float64 `json:"amount"`
	Rank   int     `json:"rank"`
}

// MarqueeItem 中奖跑马灯(全站真实记录, 昵称已脱敏)。
type MarqueeItem struct {
	Nickname  string  `json:"nickname"`
	PrizeName string  `json:"prize_name"`
	PrizeType int     `json:"prize_type"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

// InfoReq 活动信息(公开, 挂 AuthOptional: 带 token 额外返回我的次数与余额)。
type InfoReq struct {
	g.Meta      `path:"/lottery/info" method:"get" tags:"Front/Lottery" summary:"抽奖活动信息"`
	LotteryType int `json:"lottery_type" v:"required|in:1,2#玩法类型必填"`
}
type InfoRes struct {
	ActivityId  int64         `json:"activity_id"`
	Name        string        `json:"name"`
	LotteryType int           `json:"lottery_type"`
	PayType     int           `json:"pay_type"` // 1仅免费次数 2金币
	CostGold    float64       `json:"cost_gold"`
	DailyFree   int           `json:"daily_free"`
	DailyLimit  int           `json:"daily_limit"`
	Notice      string        `json:"notice"`
	Prizes      []PrizeItem   `json:"prizes"`
	FreeLeft    int           `json:"free_left"` // 今日剩余免费次数
	DrawLeft    int           `json:"draw_left"` // 今日剩余总次数, -1=不限
	Balance     float64       `json:"balance"`
	LoggedIn    bool          `json:"logged_in"`
	Marquee     []MarqueeItem `json:"marquee"`
}

// DrawReq 抽奖(需登录)。pay_type: 1用免费次数 2花金币。
type DrawReq struct {
	g.Meta      `path:"/lottery/draw" method:"post" tags:"Front/Lottery" summary:"抽奖"`
	LotteryType int `json:"lottery_type" v:"required|in:1,2#玩法类型必填"`
	PayType     int `json:"pay_type"     v:"required|in:1,2#消耗方式非法"`
}
type DrawRes struct {
	HistoryId   int64   `json:"history_id"`
	PrizeId     int64   `json:"prize_id"`
	PrizeName   string  `json:"prize_name"`
	PrizeType   int     `json:"prize_type"`
	PrizeAmount float64 `json:"prize_amount"`
	PrizeCover  string  `json:"prize_cover"`
	PrizeDesc   string  `json:"prize_desc"`
	Status      int     `json:"status"`    // 1已发放 2待发货
	NeedAddr    bool    `json:"need_addr"` // true=实物奖, 引导去填收货地址
	CostGold    float64 `json:"cost_gold"` // 本次实付
	FreeLeft    int     `json:"free_left"`
	DrawLeft    int     `json:"draw_left"`
	Balance     float64 `json:"balance"`
}

type MyItem struct {
	Id          int64   `json:"id"`
	PrizeName   string  `json:"prize_name"`
	PrizeType   int     `json:"prize_type"`
	PrizeText   string  `json:"prize_text"`
	PrizeAmount float64 `json:"prize_amount"`
	PayType     int     `json:"pay_type"`
	CostGold    float64 `json:"cost_gold"`
	Status      int     `json:"status"` // 1已发放 2待发货 3已发货
	Remark      string  `json:"remark"`
	CreatedAt   string  `json:"created_at"`
}

// MyReq 我的中奖记录(需登录)。lottery_type=0 表示全部。
type MyReq struct {
	g.Meta      `path:"/lottery/my" method:"get" tags:"Front/Lottery" summary:"我的中奖记录"`
	LotteryType int `json:"lottery_type"`
	Page        int `json:"page"`
	Size        int `json:"size"`
}
type MyRes struct {
	List  []MyItem `json:"list"`
	Total int      `json:"total"`
}

// AddressReq 填写实物奖收货地址(需登录, 只能填自己的实物奖且只能填一次)。
type AddressReq struct {
	g.Meta    `path:"/lottery/address" method:"post" tags:"Front/Lottery" summary:"填写收货地址"`
	HistoryId int64  `json:"history_id" v:"required|min:1#中奖记录ID必填"`
	Receiver  string `json:"receiver"   v:"required|length:1,32#收货人必填"`
	Phone     string `json:"phone"      v:"required|length:5,20#手机号必填"`
	Address   string `json:"address"    v:"required|length:5,200#收货地址必填"`
}
type AddressRes struct{}

type AddrItem struct {
	Id             int64  `json:"id"`
	HistoryId      int64  `json:"history_id"`
	PrizeName      string `json:"prize_name"`
	Receiver       string `json:"receiver"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	DeliveryStatus int    `json:"delivery_status"` // 0待填写 1待发货 2已发货
	ExpressNo      string `json:"express_no"`
	CreatedAt      string `json:"created_at"`
}

// MyAddrReq 我的收货单(需登录)。
type MyAddrReq struct {
	g.Meta `path:"/lottery/address/my" method:"get" tags:"Front/Lottery" summary:"我的收货单"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type MyAddrRes struct {
	List  []AddrItem `json:"list"`
	Total int        `json:"total"`
}
