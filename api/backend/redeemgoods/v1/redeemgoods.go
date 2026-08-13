// Package v1 后台商品兑换契约(商品管理 + 兑换记录)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Cover     string  `json:"cover"`
	Intro     string  `json:"intro"`
	CostGold  float64 `json:"cost_gold"`
	Stock     int     `json:"stock"`
	Exchanged int     `json:"exchanged"`
	Rank      int     `json:"rank"`
	Status    int     `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/redeem-goods" method:"get" tags:"Backend/Redeem" summary:"商品列表"`
	Status  string `json:"status"` // 空=全部
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta   `path:"/redeem-goods" method:"post" tags:"Backend/Redeem" summary:"新增商品"`
	Name     string  `json:"name" v:"required#商品名必填"`
	Cover    string  `json:"cover"`
	Intro    string  `json:"intro"`
	CostGold float64 `json:"cost_gold" v:"required|min:0.01#金币价必填|金币价需大于0"`
	Stock    int     `json:"stock"` // -1=不限量
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/redeem-goods/{id}" method:"put" tags:"Backend/Redeem" summary:"更新商品"`
	Id       int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name     string  `json:"name"`
	Cover    string  `json:"cover"`
	Intro    string  `json:"intro"`
	CostGold float64 `json:"cost_gold"`
	Stock    int     `json:"stock"`
	Rank     int     `json:"rank"`
	Status   int     `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/redeem-goods/{id}" method:"delete" tags:"Backend/Redeem" summary:"删除商品"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

type OrderItem struct {
	Id        int64   `json:"id"`
	UserId    int64   `json:"user_id"`
	GoodsId   int64   `json:"goods_id"`
	GoodsName string  `json:"goods_name"`
	CostGold  float64 `json:"cost_gold"`
	CreatedAt string  `json:"created_at"`
}

type OrdersReq struct {
	g.Meta  `path:"/redeem-goods/orders" method:"get" tags:"Backend/Redeem" summary:"兑换记录"`
	UserId  int64 `json:"user_id"`  // 0=全部
	GoodsId int64 `json:"goods_id"` // 0=全部
	Page    int   `json:"page"`
	Size    int   `json:"size"`
}
type OrdersRes struct {
	List  []OrderItem `json:"list"`
	Total int         `json:"total"`
}
