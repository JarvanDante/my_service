// Package v1 前台钱包契约(余额 + 流水)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// BalanceReq 我的余额(需登录)。
type BalanceReq struct {
	g.Meta `path:"/wallet/balance" method:"get" tags:"Front/Wallet" summary:"我的余额"`
}
type BalanceRes struct {
	Balance   float64 `json:"balance"`   // 可用余额
	Frozen    float64 `json:"frozen"`    // 提现在途冻结金额(申请中+审核通过未打款)
	TotalIn   float64 `json:"total_in"`  // 累计收入
	TotalOut  float64 `json:"total_out"` // 累计支出
	Withdrawn float64 `json:"withdrawn"` // 累计提现成功金额
}

type WaterItem struct {
	Id            int64   `json:"id"`
	Direction     int     `json:"direction"` // 1收入 2支出
	Scene         string  `json:"scene"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	Remark        string  `json:"remark"`
	CreatedAt     string  `json:"created_at"`
}

// WatersReq 我的流水(需登录)。direction/scene 为空表示不筛选。
type WatersReq struct {
	g.Meta    `path:"/wallet/waters" method:"get" tags:"Front/Wallet" summary:"我的金币流水"`
	Direction string `json:"direction"` // ""=全部 "1"=收入 "2"=支出
	Scene     string `json:"scene"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type WatersRes struct {
	List  []WaterItem `json:"list"`
	Total int         `json:"total"`
}
