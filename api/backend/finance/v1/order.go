// Package v1 后台财务接口契约(B2): 订单/流水/支付回调。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 充值订单 ----------

type OrderItem struct {
	Id        int64   `json:"id"`
	OrderNo   string  `json:"order_no"`
	UserId    int64   `json:"user_id"`
	PackageId int64   `json:"package_id"`
	Amount    float64 `json:"amount"`
	Coin      float64 `json:"coin"`
	Status    int     `json:"status"` // 0待支付 1已支付 -1取消
	PayAt     string  `json:"pay_at"`
	CreatedAt string  `json:"created_at"`
}

type OrderListReq struct {
	g.Meta    `path:"/recharge-orders" method:"get" tags:"Backend/Finance" summary:"充值订单列表"`
	OrderNo   string `json:"order_no"`
	UserId    int64  `json:"user_id"    v:"min:0#用户ID不合法"`
	Status    int    `json:"status"     v:"in:0,1,2,3#状态仅支持0/1/2/3"` // 0全部 1待支付 2已支付 3已取消
	StartDate string `json:"start_date"`                             // YYYY-MM-DD
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type OrderListRes struct {
	List  []OrderItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// ---------- 全站余额流水 ----------

type BalanceLogItem struct {
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

type BalanceLogListReq struct {
	g.Meta    `path:"/balance-logs" method:"get" tags:"Backend/Finance" summary:"全站余额流水"`
	UserId    int64  `json:"user_id"   v:"min:0#用户ID不合法"`
	Scene     string `json:"scene"`
	Direction int    `json:"direction" v:"in:0,1,2#direction 仅支持0/1/2"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type BalanceLogListRes struct {
	List  []BalanceLogItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// ---------- 支付回调(公开, 签名校验) ----------

type PayCallbackReq struct {
	g.Meta  `path:"/pay/callback" method:"post" tags:"Backend/Finance" summary:"支付回调(网关调用)"`
	OrderNo string `json:"order_no" v:"required#订单号必填"`
	Sign    string `json:"sign"` // md5(order_no + secret); 未配置 secret 时(开发环境)可不传
}
type PayCallbackRes struct{}
