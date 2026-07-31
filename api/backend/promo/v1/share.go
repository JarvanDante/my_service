// Package v1 后台分享/拉新接口契约(B6)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 分享记录 ----------

type ShareLogAdminItem struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	Type      string `json:"type"`
	TargetId  int64  `json:"target_id"`
	Channel   string `json:"channel"`
	CreatedAt string `json:"created_at"`
}

type ShareLogListReq struct {
	g.Meta    `path:"/share-logs" method:"get" tags:"Backend/Promo" summary:"分享记录"`
	UserId    int64  `json:"user_id" v:"min:0#user_id 不合法"`
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type ShareLogListRes struct {
	List  []ShareLogAdminItem `json:"list"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

// ---------- 分享统计 / 拉新排行 ----------

type ChannelCount struct {
	Channel string `json:"channel"`
	Count   int    `json:"count"`
}

type InviteRankItem struct {
	UserId      int64  `json:"user_id"`
	Username    string `json:"username"`
	InviteCount int    `json:"invite_count"`
}

type ShareStatsReq struct {
	g.Meta    `path:"/share-stats" method:"get" tags:"Backend/Promo" summary:"分享统计与拉新排行"`
	StartDate string `json:"start_date"` // YYYY-MM-DD, 可选
	EndDate   string `json:"end_date"`
	Top       int    `json:"top" v:"min:0#top 不合法"` // 拉新排行条数, 默认10, 最大100
}
type ShareStatsRes struct {
	TotalShares int              `json:"total_shares"` // 分享总次数
	SharerCount int              `json:"sharer_count"` // 分享人数
	Channels    []ChannelCount   `json:"channels"`     // 渠道分布
	InviteRank  []InviteRankItem `json:"invite_rank"`  // 拉新排行
}
