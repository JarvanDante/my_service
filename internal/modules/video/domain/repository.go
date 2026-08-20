package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type ListFilter struct {
	Keyword   string
	MediaCode string
	Kind      int // 0视频 2动漫
	Status    int // 9=全部
	// Sort 排序方式: 0综合(人工 sort 权重, 后台与前台默认) 1最新 2时长。
	// 加在仓储层而不是在 logic 里排内存, 是因为分页必须由 SQL 完成 ——
	// 内存排序只能排当前这一页, 结果是错的。
	Sort int
}

type Repository interface {
	List(ctx context.Context, f ListFilter, page, size int) ([]*entity.Video, int, error)
	Find(ctx context.Context, id int64) (*entity.Video, error)
	FindByMediaCode(ctx context.Context, code string, kind int) (*entity.Video, error)
	Create(ctx context.Context, v *entity.Video) (int64, error)
	Update(ctx context.Context, v *entity.Video) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status int) error
}
