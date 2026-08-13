// Package v1 前台优惠券契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type TplItem struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`  // 1抵用券 2折扣券
	Scene     int     `json:"scene"` // 1充值 2内容购买 3通用
	FaceValue float64 `json:"face_value"`
	Discount  int     `json:"discount"`
	Threshold float64 `json:"threshold"`
	MaxDeduct float64 `json:"max_deduct"`
	ExpireDay int     `json:"expire_day"`
	Received  bool    `json:"received"` // 当前用户是否已领到上限
}

// TplsReq 可领取的券模板(公开)。
type TplsReq struct {
	g.Meta `path:"/coupon/tpls" method:"get" tags:"Front/Coupon" summary:"可领券列表"`
}
type TplsRes struct {
	List []TplItem `json:"list"`
}

// ReceiveReq 领券(需登录, 受每人限领与总量约束)。
type ReceiveReq struct {
	g.Meta `path:"/coupon/receive" method:"post" tags:"Front/Coupon" summary:"领取优惠券"`
	TplId  int64 `json:"tpl_id" v:"required|min:1#券模板ID必填"`
}
type ReceiveRes struct {
	Id int64 `json:"id"`
}

type MyItem struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Type       int     `json:"type"`
	Scene      int     `json:"scene"`
	FaceValue  float64 `json:"face_value"`
	Discount   int     `json:"discount"`
	Threshold  float64 `json:"threshold"`
	MaxDeduct  float64 `json:"max_deduct"`
	Status     int     `json:"status"` // 1未使用 2已使用 3已过期
	StatusText string  `json:"status_text"`
	ExpireAt   string  `json:"expire_at"`
	UsedAt     string  `json:"used_at"`
	CreatedAt  string  `json:"created_at"`
}

// MyReq 我的券(需登录)。status 空=全部。
type MyReq struct {
	g.Meta `path:"/coupon/my" method:"get" tags:"Front/Coupon" summary:"我的优惠券"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type MyRes struct {
	List  []MyItem `json:"list"`
	Total int      `json:"total"`
}

// BestReq 给定订单金额与场景, 返回可用券列表与最优券(服务端算抵扣, 客户端不参与)。
type BestReq struct {
	g.Meta `path:"/coupon/available" method:"get" tags:"Front/Coupon" summary:"可用券与最优券"`
	Scene  int     `json:"scene" v:"in:0,1,2,3#场景非法"`
	Amount float64 `json:"amount" v:"required|min:0.01#订单金额必填"`
}
type AvailableItem struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Deduct float64 `json:"deduct"` // 该券对本单的抵扣额
	Type   int     `json:"type"`
	Expire string  `json:"expire_at"`
	IsBest bool    `json:"is_best"`
}
type BestRes struct {
	List       []AvailableItem `json:"list"`
	BestId     int64           `json:"best_id"`
	BestDeduct float64         `json:"best_deduct"`
	PayAmount  float64         `json:"pay_amount"` // 用最优券后的应付金额
}
