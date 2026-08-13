// Package v1 前台图集契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Pic struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Item 列表项。列表不下发 pics(图片体积大且列表页用不到), 只给封面与总张数。
type Item struct {
	Id        int64    `json:"id"`
	Title     string   `json:"title"`
	Cover     string   `json:"cover"`
	Intro     string   `json:"intro"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	IsVip     int      `json:"is_vip"`
	Price     float64  `json:"price"`
	FreeCount int      `json:"free_count"`
	PicCount  int      `json:"pic_count"`
	ViewCount int64    `json:"view_count"`
	LikeCount int64    `json:"like_count"`
	IsBuy     bool     `json:"is_buy"`
	CreatedAt string   `json:"created_at"`
}

// ListReq 图集列表(公开)。sort: 0综合(rank) 1最多观看 2最新 3最多点赞。
type ListReq struct {
	g.Meta   `path:"/photo/list" method:"get" tags:"Front/Photo" summary:"图集列表"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
	Keyword  string `json:"keyword"`
	Sort     int    `json:"sort"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// DetailReq 图集详情(公开; 带登录态则返回解锁信息)。
type DetailReq struct {
	g.Meta `path:"/photo/detail" method:"get" tags:"Front/Photo" summary:"图集详情"`
	Id     int64 `json:"id" v:"required|min:1#图集ID必填"`
}

// DetailRes 未解锁时 pics 只有前 free_count 张(服务端截断),
// preview_count/total_count 让前端明确知道"这是被截断的预览", 好画"解锁查看剩余 N 张"。
type DetailRes struct {
	Item
	Pics         []Pic  `json:"pics"`
	Playable     bool   `json:"playable"` // 整套是否已解锁
	NeedPay      bool   `json:"need_pay"`
	NeedVip      bool   `json:"need_vip"`
	Enough       bool   `json:"enough"` // 余额是否够买
	Reason       string `json:"reason"`
	PreviewCount int    `json:"preview_count"` // 本次实际下发的张数
	TotalCount   int    `json:"total_count"`   // 图集总张数
}

// BuyReq 整套购买(需登录)。价格服务端定, 不接受客户端传金额。
type BuyReq struct {
	g.Meta `path:"/photo/buy" method:"post" tags:"Front/Photo" summary:"购买图集"`
	Id     int64 `json:"id" v:"required|min:1#图集ID必填"`
}
type BuyRes struct {
	Price   float64 `json:"price"`
	Balance float64 `json:"balance"` // 购买后余额
}
