// Package service 排行/热搜对外接口。
package service

import "context"

type RankItem struct {
	ContentId int64
	MediaType int
	Score     int64
	RankNo    int
}

type HotDTO struct {
	Id          int64
	Keyword     string
	Heat        int
	SearchCount int64
	Status      int
	UpdatedAt   string
}

type HotFilter struct {
	Status  int // -1=全部
	Keyword string
	Page    int
	Size    int
}

type IRank interface {
	// Rank 点赞聚合排行(period=day/week/all), Redis 缓存 60s, 失败降级直查。
	Rank(ctx context.Context, mediaType int, period string) ([]RankItem, error)
	// RefreshRank 清缓存。
	RefreshRank(ctx context.Context) error
	// HotKeywords 前台热搜词(heat desc, search_count desc, 上限 20)。
	HotKeywords(ctx context.Context) ([]string, error)
	// 后台热搜词管理
	HotList(ctx context.Context, f HotFilter) ([]*HotDTO, int, error)
	HotCreate(ctx context.Context, keyword string, heat, status int) (int64, error)
	HotUpdate(ctx context.Context, id int64, keyword string, heat, status int) error
	HotDelete(ctx context.Context, id int64) error
}
