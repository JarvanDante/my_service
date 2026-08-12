// Package v1 前台兑换码契约(移植自 tianbi redeemcode use/list)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// UseReq 使用兑换码(需登录)。
type UseReq struct {
	g.Meta `path:"/redeemcode/use" method:"post" tags:"Front/RedeemCode" summary:"使用兑换码"`
	Code   string `json:"code" v:"required#兑换码必填"`
}
type UseRes struct {
	Desc string `json:"desc"` // 如: 兑换金币x100
}

// RecordItem 我的兑换记录(与 tianbi AppRedeemCodeRecord 对齐)。
type RecordItem struct {
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	ActivedAt string `json:"actived_at"`
}

// ListReq 我的兑换记录(需登录)。
type ListReq struct {
	g.Meta `path:"/redeemcode/list" method:"get" tags:"Front/RedeemCode" summary:"我的兑换记录"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type ListRes struct {
	List []RecordItem `json:"list"`
}
