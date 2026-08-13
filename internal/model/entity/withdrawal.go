// Code maintained manually (提现单 + 用户收款账户).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// Withdrawal 提现单。状态: 1申请中 2审核通过 4已打款 5已拒绝 6已撤回。
type Withdrawal struct {
	Id           int64       `json:"id"           orm:"id"`
	SiteId       int64       `json:"siteId"       orm:"site_id"`
	TradeNo      string      `json:"tradeNo"      orm:"trade_no"`
	UserId       int64       `json:"userId"       orm:"user_id"`
	Amount       float64     `json:"amount"       orm:"amount"`
	Fee          float64     `json:"fee"          orm:"fee"`
	RealAmount   float64     `json:"realAmount"   orm:"real_amount"`
	FeeRate      float64     `json:"feeRate"      orm:"fee_rate"`
	BalanceAfter float64     `json:"balanceAfter" orm:"balance_after"`
	AccountInfo  string      `json:"accountInfo"  orm:"account_info"` // jsonb 原文
	Status       int         `json:"status"       orm:"status"`
	AuditBy      int64       `json:"auditBy"      orm:"audit_by"`
	AuditAt      *gtime.Time `json:"auditAt"      orm:"audit_at"`
	PaidAt       *gtime.Time `json:"paidAt"       orm:"paid_at"`
	Remark       string      `json:"remark"       orm:"remark"`
	PayVoucher   string      `json:"payVoucher"   orm:"pay_voucher"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}

// BankCard 用户收款账户。account_type: 1银行卡 2支付宝 3USDT。
type BankCard struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	UserId      int64       `json:"userId"      orm:"user_id"`
	AccountType int         `json:"accountType" orm:"account_type"`
	AccountName string      `json:"accountName" orm:"account_name"`
	AccountNo   string      `json:"accountNo"   orm:"account_no"`
	BankName    string      `json:"bankName"    orm:"bank_name"`
	IsDefault   int         `json:"isDefault"   orm:"is_default"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
