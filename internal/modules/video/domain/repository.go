package domain

import (
	"context"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type ListFilter struct {
	Keyword      string
	MediaCode    string
	Category     string   // 作品分类名, 命中 category 逗号列表中的一项
	Categories   []string // 多个分类名, 命中任一项
	Kind         int      // 0视频 2动漫 3抖音
	Status       int      // 9=全部 8=草稿+下架
	SubmitSource int      // 9=全部 0后台 1用户
	// Sort 排序方式: 0综合(人工 sort 权重, 后台与前台默认) 1最新 2时长。
	// 加在仓储层而不是在 logic 里排内存, 是因为分页必须由 SQL 完成 ——
	// 内存排序只能排当前这一页, 结果是错的。
	Sort      int
	Tags      []string // 标签名, 命中 tags jsonb 数组中任一项
	UpUserIds []int64  // 只出这些 UP 主的作品(关注流)
}

type Repository interface {
	List(ctx context.Context, f ListFilter, page, size int) ([]*entity.Video, int, error)
	Find(ctx context.Context, id int64) (*entity.Video, error)
	FindByMediaCode(ctx context.Context, code string, kind int) (*entity.Video, error)
	Create(ctx context.Context, v *entity.Video) (int64, error)
	Update(ctx context.Context, v *entity.Video) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status int) error
	Audit(ctx context.Context, id int64, status int, reason string, operatorId int64) (int64, error)
}
