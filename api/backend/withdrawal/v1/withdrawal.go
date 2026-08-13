// Package v1 后台提现契约(列表 + 审核 + 打款)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64   `json:"id"`
	TradeNo      string  `json:"trade_no"`
	UserId       int64   `json:"user_id"`
	Nickname     string  `json:"nickname"`
	Amount       float64 `json:"amount"`
	Fee          float64 `json:"fee"`
	RealAmount   float64 `json:"real_amount"`
	FeeRate      float64 `json:"fee_rate"`
	BalanceAfter float64 `json:"balance_after"`
	AccountName  string  `json:"account_name"`
	AccountNo    string  `json:"account_no"`
	BankName     string  `json:"bank_name"`
	Status       int     `json:"status"`
	StatusText   string  `json:"status_text"`
	AuditBy      int64   `json:"audit_by"`
	AuditAt      string  `json:"audit_at"`
	PaidAt       string  `json:"paid_at"`
	Remark       string  `json:"remark"`
	PayVoucher   string  `json:"pay_voucher"`
	CreatedAt    string  `json:"created_at"`
}

// ListReq 提现单列表。status/user_id 用 string 接收, 空=全部。
type ListReq struct {
	g.Meta `path:"/withdrawals" method:"get" tags:"Backend/Withdrawal" summary:"提现单列表"`
	Status string `json:"status"`
	UserId string `json:"user_id"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type ListRes struct {
	List       []Item  `json:"list"`
	Total      int     `json:"total"`
	SumAmount  float64 `json:"sum_amount"`  // 当前筛选条件下的申请总额
	PendingNum int     `json:"pending_num"` // 待审核笔数(status=1)
}

// AuditReq 审核。pass=true → 2审核通过; pass=false → 5拒绝(自动退款)。
// 仅 status=1 可审核, 条件更新保证并发下只生效一次。
type AuditReq struct {
	g.Meta `path:"/withdrawals/{id}/audit" method:"post" tags:"Backend/Withdrawal" summary:"审核提现"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Pass   bool   `json:"pass"`
	Remark string `json:"remark"`
}
type AuditRes struct{}

// MarkPaidReq 标记已打款(仅 status=2 可打款, 终态)。
type MarkPaidReq struct {
	g.Meta  `path:"/withdrawals/{id}/mark-paid" method:"post" tags:"Backend/Withdrawal" summary:"标记已打款"`
	Id      int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Voucher string `json:"voucher"`
	Remark  string `json:"remark"`
}
type MarkPaidRes struct{}

// RejectPaidReq 打款失败退款(status=2 → 5拒绝并退款), 用于线下打款失败的兜底。
type RejectPaidReq struct {
	g.Meta `path:"/withdrawals/{id}/refund" method:"post" tags:"Backend/Withdrawal" summary:"打款失败退款"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Remark string `json:"remark"`
}
type RejectPaidRes struct{}
