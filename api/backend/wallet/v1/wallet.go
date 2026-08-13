// Package v1 后台钱包契约(全站流水 + 人工调账)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type LogItem struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	Direction     int     `json:"direction"`
	Scene         string  `json:"scene"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	RefId         string  `json:"ref_id"`
	Remark        string  `json:"remark"`
	CreatedAt     string  `json:"created_at"`
}

// LogsReq 全站金币流水。user_id/direction 用 string 接收, 空=全部(避免 0 与空歧义)。
type LogsReq struct {
	g.Meta    `path:"/wallet/logs" method:"get" tags:"Backend/Wallet" summary:"金币流水"`
	UserId    string `json:"user_id"`
	Direction string `json:"direction"`
	Scene     string `json:"scene"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type LogsRes struct {
	List  []LogItem `json:"list"`
	Total int       `json:"total"`
}

// AdjustReq 人工调账: amount 正数=加币, 负数=扣币(扣币受余额约束, 不允许扣成负数)。
type AdjustReq struct {
	g.Meta `path:"/wallet/adjust" method:"post" tags:"Backend/Wallet" summary:"人工调账"`
	UserId int64   `json:"user_id" v:"required|min:1#用户ID必填"`
	Amount float64 `json:"amount"  v:"required#调账金额必填"`
	Remark string  `json:"remark"  v:"required#调账备注必填"`
}
type AdjustRes struct {
	Balance float64 `json:"balance"` // 调账后余额
}
