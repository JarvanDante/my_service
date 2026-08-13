// Code maintained manually (兑换商品 + 兑换记录).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type RedeemGoods struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	Cover     string      `json:"cover"     orm:"cover"`
	Intro     string      `json:"intro"     orm:"intro"`
	CostGold  float64     `json:"costGold"  orm:"cost_gold"`
	Stock     int         `json:"stock"     orm:"stock"` // -1=不限量
	Exchanged int         `json:"exchanged" orm:"exchanged"`
	Rank      int         `json:"rank"      orm:"rank"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

type RedeemGoodsOrder struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	GoodsId   int64       `json:"goodsId"   orm:"goods_id"`
	GoodsName string      `json:"goodsName" orm:"goods_name"`
	CostGold  float64     `json:"costGold"  orm:"cost_gold"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
