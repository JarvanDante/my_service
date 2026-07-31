// Package domain 系统模块领域层(B7: 公告/推送)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

// NoticeFilter 公告筛选。
type NoticeFilter struct {
	Type   string // notice / push
	Status int    // 0全部 1上架 2下线
}

type Repository interface {
	NoticeCreate(ctx context.Context, n *entity.SystemNotice) (int64, error)
	NoticeList(ctx context.Context, f NoticeFilter, page, size int) ([]*entity.SystemNotice, int, error)
	NoticeFind(ctx context.Context, id int64) (*entity.SystemNotice, error)
	NoticeSetStatus(ctx context.Context, id int64, status int) error
}
