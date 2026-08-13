// Package service 提现对外接口。
package service

import "context"

// 提现状态机(与迁移 00033 注释一致)
const (
	StatusApplying = 1 // 申请中(已冻结)
	StatusPassed   = 2 // 审核通过, 待打款
	StatusPaid     = 4 // 已打款(终态)
	StatusRejected = 5 // 已拒绝(已退款, 终态)
	StatusCanceled = 6 // 用户撤回(已退款, 终态)
)

// StatusText 状态中文, 前后台共用。
func StatusText(s int) string {
	switch s {
	case StatusApplying:
		return "审核中"
	case StatusPassed:
		return "待打款"
	case StatusPaid:
		return "已打款"
	case StatusRejected:
		return "已拒绝"
	case StatusCanceled:
		return "已撤回"
	}
	return "未知"
}

type ConfigDTO struct {
	Open       bool
	MinAmount  float64
	MaxAmount  float64
	Multiple   float64
	FeeRate    float64
	DailyLimit int
	DailyUsed  int
	Balance    float64
	Frozen     float64
}

type OrderDTO struct {
	Id           int64
	TradeNo      string
	UserId       int64
	Nickname     string
	Amount       float64
	Fee          float64
	RealAmount   float64
	FeeRate      float64
	BalanceAfter float64
	AccountName  string
	AccountNo    string
	BankName     string
	Status       int
	AuditBy      int64
	AuditAt      string
	PaidAt       string
	Remark       string
	PayVoucher   string
	CreatedAt    string
}

type ApplyResult struct {
	Id         int64
	TradeNo    string
	Amount     float64
	Fee        float64
	RealAmount float64
}

type ListFilter struct {
	UserId int64 // 0=全部
	Status int   // -1=全部
	Page   int
	Size   int
}

type CardDTO struct {
	Id          int64
	AccountType int
	AccountName string
	AccountNo   string
	BankName    string
	IsDefault   int
	CreatedAt   string
}

type CardInput struct {
	Id          int64
	UserId      int64
	AccountType int
	AccountName string
	AccountNo   string
	BankName    string
	IsDefault   int
}

type IWithdrawal interface {
	Config(ctx context.Context, userId int64) (*ConfigDTO, error)
	// Apply 申请提现: 校验(开关/额度/倍数/频率/日限) → 事务(条件扣款冻结 + 写流水 + 建单)。
	Apply(ctx context.Context, userId, cardId int64, amount float64) (*ApplyResult, error)
	My(ctx context.Context, f ListFilter) ([]*OrderDTO, int, error)
	// Cancel 用户撤回: 仅 status=1, 条件更新 + 原路退款。
	Cancel(ctx context.Context, userId, id int64) error

	List(ctx context.Context, f ListFilter) ([]*OrderDTO, int, float64, int, error)
	// Audit 审核: pass → 2; 否则 → 5 并退款。仅 status=1 生效。
	Audit(ctx context.Context, adminId, id int64, pass bool, remark string) error
	// MarkPaid 标记打款: 仅 status=2 生效, 不退款。
	MarkPaid(ctx context.Context, adminId, id int64, voucher, remark string) error
	// RefundPaid 打款失败退款: 仅 status=2 生效, 退款并置 5。
	RefundPaid(ctx context.Context, adminId, id int64, remark string) error

	CardList(ctx context.Context, userId int64) ([]*CardDTO, error)
	CardAdd(ctx context.Context, in CardInput) (int64, error)
	CardUpdate(ctx context.Context, in CardInput) error
	CardDel(ctx context.Context, userId int64, ids []int64) error
}
