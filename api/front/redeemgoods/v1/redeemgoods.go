// Package v1 前台商品兑换契约(移植自 tianbi redeem)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type GoodsItem struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Cover    string  `json:"cover"`
	Intro    string  `json:"intro"`
	CostGold float64 `json:"cost_gold"`
	Stock    int     `json:"stock"` // -1=不限量
}

// ListReq 可兑换商品列表(公开)。
type ListReq struct {
	g.Meta `path:"/redeem/list" method:"get" tags:"Front/Redeem" summary:"兑换商品列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type ListRes struct {
	List  []GoodsItem `json:"list"`
	Total int         `json:"total"`
}

// ExchangeReq 兑换(需登录, 事务: 扣金币+减库存+记录)。
type ExchangeReq struct {
	g.Meta  `path:"/redeem/exchange" method:"post" tags:"Front/Redeem" summary:"兑换商品"`
	GoodsId int64 `json:"goods_id" v:"required|min:1#商品ID必填"`
}
type ExchangeRes struct {
	OrderId int64 `json:"order_id"`
}

type HistoryItem struct {
	Id        int64   `json:"id"`
	GoodsName string  `json:"goods_name"`
	CostGold  float64 `json:"cost_gold"`
	CreatedAt string  `json:"created_at"`
}

// HistoryReq 我的兑换历史(需登录)。
type HistoryReq struct {
	g.Meta `path:"/redeem/history" method:"get" tags:"Front/Redeem" summary:"我的兑换历史"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type HistoryRes struct {
	List []HistoryItem `json:"list"`
}
