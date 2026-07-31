// Package domain 推广/兑换码模块领域层(B3)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

// CodeFilter 兑换码筛选。
type CodeFilter struct {
	Keyword string // 码/名称 模糊
	CodeKey string // 批次
	Type    string // point / group
	Status  int    // 0全部 1可用 2已使用 3作废
}

// CodeLogFilter 兑换记录筛选。
type CodeLogFilter struct {
	CodeId int64
	UserId int64
	Code   string
}

// ShareLogFilter 分享记录筛选(B6)。
type ShareLogFilter struct {
	UserId    int64
	Type      string
	Channel   string
	StartDate string // created_at >= (YYYY-MM-DD)
	EndDate   string // created_at < 次日
}

// ChannelCount 渠道分布(B6)。
type ChannelCount struct {
	Channel string
	Count   int
}

// InviteRankItem 拉新排行项(B6)。
type InviteRankItem struct {
	UserId      int64
	Username    string
	InviteCount int
}

type Repository interface {
	ListCodes(ctx context.Context, f CodeFilter, page, size int) ([]*entity.UserCode, int, error)
	FindCodeById(ctx context.Context, id int64) (*entity.UserCode, error)
	BatchCreateCodes(ctx context.Context, rows []*entity.UserCode) error
	VoidCode(ctx context.Context, id int64) error
	ListCodeLogs(ctx context.Context, f CodeLogFilter, page, size int) ([]*entity.UserCodeLog, int, error)

	// 分享 / 拉新(B6)
	ShareLogList(ctx context.Context, f ShareLogFilter, page, size int) ([]*entity.UserShareLog, int, error)
	ShareStats(ctx context.Context, startDate, endDate string) (totalShares int, sharerCount int, channels []ChannelCount, err error)
	InviteRank(ctx context.Context, startDate, endDate string, top int) ([]*InviteRankItem, error)
}
