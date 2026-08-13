// Package v1 前台全站搜索契约。
//
// 一个接口同时承担两种形态(type=0 聚合 / type>0 单类分页), 所以响应结构里
// 两组字段共存: 分组切片走 omitempty, 用不到的那一支直接从 JSON 里消失,
// 前端不会拿到一堆 null 需要判空; 计数字段不加 omitempty —— 搜不到时
// total/total_hit 必须真实出现且为 0, 否则前端分不清"没结果"和"字段没返回"。
package v1

import "github.com/gogf/gf/v2/frame/g"

// Item 跨表统一的结果项。
// 各内容表字段不齐(video 没有 author/price/is_vip, post/photo 没有 author),
// 缺的一律给零值, 让前端只认这一套结构, 不用按 media_type 切换解析逻辑。
type Item struct {
	Id        int64   `json:"id"`
	MediaType int     `json:"media_type"` // 1视频 2帖子 3漫画 4小说 5图集(与 paywall 同一套编码)
	Title     string  `json:"title"`
	Cover     string  `json:"cover"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	IsVip     int     `json:"is_vip"`
	ViewCount int64   `json:"view_count"`
	CreatedAt string  `json:"created_at"`
}

// SearchReq 全站搜索(公开, 挂 AuthOptional)。
// type=0 时每类各取一小撮聚合返回; type>0 时只搜该类并按 page/size 分页。
type SearchReq struct {
	g.Meta  `path:"/search" method:"get" tags:"Front/Search" summary:"全站搜索"`
	Keyword string `json:"keyword" v:"required#关键词必填"`
	Type    int    `json:"type" v:"in:0,1,2,3,4,5#资源类型非法" dc:"0全部 1视频 2帖子 3漫画 4小说 5图集"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type SearchRes struct {
	// ---- type=0 聚合分组(type>0 时这几个字段不下发)----
	Videos []Item `json:"videos,omitempty"`
	Posts  []Item `json:"posts,omitempty"`
	Comics []Item `json:"comics,omitempty"`
	Novels []Item `json:"novels,omitempty"`
	Photos []Item `json:"photos,omitempty"`
	// ---- type>0 单类分页(type=0 时不下发)----
	List []Item `json:"list,omitempty"`
	// TotalHit 全类命中总数(type=0 用), Total 当前类命中总数(type>0 用)。
	TotalHit int `json:"total_hit"`
	Total    int `json:"total"`
}

// SuggestReq 搜索联想(公开)。从热搜词表按前缀匹配, 不查内容表 ——
// 联想要的是"别人搜过什么", 内容标题前缀匹配反而噪音大且没法排序。
type SuggestReq struct {
	g.Meta  `path:"/search/suggest" method:"get" tags:"Front/Search" summary:"搜索联想词"`
	Keyword string `json:"keyword"`
}
type SuggestRes struct {
	List []string `json:"list"`
}
