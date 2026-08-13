// Package service UGC投稿对外接口。
package service

import "context"

type PublishDTO struct {
	Id           int64
	UserId       int64
	Type         int
	Title        string
	Intro        string
	Cover        string
	Resource     []string
	Tags         []string
	Status       int
	RejectReason string
	AuditBy      int64
	AuditAt      string
	CreatedAt    string
}

type SubmitInput struct {
	UserId   int64
	Type     int
	Title    string
	Intro    string
	Cover    string
	Resource []string
	Tags     []string
}

// ListFilter 前后台共用。UserId/Type=0 与 Status=-1 都表示"不筛选",
// 因为后台参数按项目铁律用 string 接收, 空串在控制器里转成这些哨兵值。
type ListFilter struct {
	UserId  int64
	Type    int
	Status  int
	Keyword string
	Page    int
	Size    int
}

type IPublish interface {
	// Submit 提交投稿: 过敏感词, 落库为待审。
	Submit(ctx context.Context, in SubmitInput) (int64, error)
	// My 我的投稿(含被拒/已撤回)。
	My(ctx context.Context, userId int64, f ListFilter) ([]*PublishDTO, int, error)
	// Cancel 撤回自己的待审投稿(条件更新, 重复撤回会失败)。
	Cancel(ctx context.Context, userId, id int64) error

	List(ctx context.Context, f ListFilter) ([]*PublishDTO, int, error)
	// Audit 审核(条件更新, 仅待审可审; adminId 留痕)。
	Audit(ctx context.Context, id, adminId int64, pass bool, rejectReason string) error
}
