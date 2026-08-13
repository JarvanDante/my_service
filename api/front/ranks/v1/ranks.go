// Package v1 前台排行/热搜契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type RankItem struct {
	ContentId int64 `json:"content_id"`
	MediaType int   `json:"media_type"`
	Score     int64 `json:"score"` // 点赞数
	RankNo    int   `json:"rank_no"`
}

// RankReq 内容排行(公开; 点赞聚合, Redis 缓存 60s)。
type RankReq struct {
	g.Meta    `path:"/rank/list" method:"get" tags:"Front/Rank" summary:"内容排行榜"`
	MediaType int    `json:"media_type" v:"required|in:1,2#资源类型必填|资源类型非法"` // 1视频 2帖子
	Period    string `json:"period"`                                       // day/week/all(默认)
}
type RankRes struct {
	List []RankItem `json:"list"`
}

type HotItem struct {
	Keyword string `json:"keyword"`
}

// HotReq 热搜词(公开)。
type HotReq struct {
	g.Meta `path:"/hotsearch/list" method:"get" tags:"Front/Rank" summary:"热搜词"`
}
type HotRes struct {
	List []HotItem `json:"list"`
}
