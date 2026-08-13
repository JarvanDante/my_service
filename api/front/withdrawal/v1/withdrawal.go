// Package v1 前台提现契约(移植自 tianbi withdrawal + bankcard)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ConfigReq 提现配置 + 我的可提余额(需登录)。
type ConfigReq struct {
	g.Meta `path:"/withdrawal/config" method:"get" tags:"Front/Withdrawal" summary:"提现配置"`
}
type ConfigRes struct {
	Open       bool    `json:"open"`        // 提现总开关
	MinAmount  float64 `json:"min_amount"`  // 单笔最低
	MaxAmount  float64 `json:"max_amount"`  // 单笔最高
	Multiple   float64 `json:"multiple"`    // 金额倍数约束, 0=不限制
	FeeRate    float64 `json:"fee_rate"`    // 手续费率(%)
	DailyLimit int     `json:"daily_limit"` // 每日次数上限
	DailyUsed  int     `json:"daily_used"`  // 今日已申请次数
	Balance    float64 `json:"balance"`     // 当前可用余额
	Frozen     float64 `json:"frozen"`      // 在途冻结
}

// ApplyReq 发起提现(需登录)。金额从余额扣除并冻结, 手续费服务端算, 不信客户端。
type ApplyReq struct {
	g.Meta `path:"/withdrawal/apply" method:"post" tags:"Front/Withdrawal" summary:"申请提现"`
	CardId int64   `json:"card_id" v:"required|min:1#收款账户必选"`
	Amount float64 `json:"amount"  v:"required|min:0.01#提现金额必填"`
}
type ApplyRes struct {
	Id         int64   `json:"id"`
	TradeNo    string  `json:"trade_no"`
	Amount     float64 `json:"amount"`
	Fee        float64 `json:"fee"`
	RealAmount float64 `json:"real_amount"`
}

type Item struct {
	Id          int64   `json:"id"`
	TradeNo     string  `json:"trade_no"`
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	RealAmount  float64 `json:"real_amount"`
	Status      int     `json:"status"` // 1申请中 2审核通过 4已打款 5已拒绝 6已撤回
	StatusText  string  `json:"status_text"`
	AccountNo   string  `json:"account_no"`
	AccountName string  `json:"account_name"`
	BankName    string  `json:"bank_name"`
	Remark      string  `json:"remark"`
	CreatedAt   string  `json:"created_at"`
}

// MyReq 我的提现记录(需登录)。status 空=全部。
type MyReq struct {
	g.Meta `path:"/withdrawal/my" method:"get" tags:"Front/Withdrawal" summary:"我的提现记录"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type MyRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// CancelReq 撤回提现(仅申请中可撤, 原路退款)。
type CancelReq struct {
	g.Meta `path:"/withdrawal/cancel" method:"post" tags:"Front/Withdrawal" summary:"撤回提现"`
	Id     int64 `json:"id" v:"required|min:1#提现单ID必填"`
}
type CancelRes struct{}

type CardItem struct {
	Id          int64  `json:"id"`
	AccountType int    `json:"account_type"` // 1银行卡 2支付宝 3USDT
	AccountName string `json:"account_name"`
	AccountNo   string `json:"account_no"`
	BankName    string `json:"bank_name"`
	IsDefault   int    `json:"is_default"`
	CreatedAt   string `json:"created_at"`
}

// CardListReq 我的收款账户(需登录)。
type CardListReq struct {
	g.Meta `path:"/bankcard/list" method:"get" tags:"Front/Withdrawal" summary:"收款账户列表"`
}
type CardListRes struct {
	List []CardItem `json:"list"`
}

// CardAddReq 添加收款账户。字段禁止含空格(移植自 tianbi 校验)。
type CardAddReq struct {
	g.Meta      `path:"/bankcard/add" method:"post" tags:"Front/Withdrawal" summary:"添加收款账户"`
	AccountType int    `json:"account_type" v:"in:0,1,2,3#账户类型非法"`
	AccountName string `json:"account_name" v:"required#开户人必填"`
	AccountNo   string `json:"account_no"   v:"required#账号必填"`
	BankName    string `json:"bank_name"`
	IsDefault   int    `json:"is_default"`
}
type CardAddRes struct {
	Id int64 `json:"id"`
}

// CardUpdateReq 修改收款账户(仅本人)。
type CardUpdateReq struct {
	g.Meta      `path:"/bankcard/update" method:"post" tags:"Front/Withdrawal" summary:"修改收款账户"`
	Id          int64  `json:"id" v:"required|min:1#ID必填"`
	AccountType int    `json:"account_type" v:"in:0,1,2,3#账户类型非法"`
	AccountName string `json:"account_name"`
	AccountNo   string `json:"account_no"`
	BankName    string `json:"bank_name"`
	IsDefault   int    `json:"is_default"`
}
type CardUpdateRes struct{}

// CardDelReq 删除收款账户(按 ids 批量, 仅本人)。
type CardDelReq struct {
	g.Meta `path:"/bankcard/del" method:"post" tags:"Front/Withdrawal" summary:"删除收款账户"`
	Ids    []int64 `json:"ids" v:"required#ID列表必填"`
}
type CardDelRes struct{}
