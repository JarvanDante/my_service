// Package service 帖子对外接口。
package service

import "context"

type PostDTO struct {
	Id           int64
	UserId       int64
	Nickname     string
	Img          string
	Title        string
	Content      string
	Pics         []string
	MediaId      int64
	ViewCount    int64
	LikeCount    int
	CommentCount int
	Status       int
	RejectReason string
	CreatedAt    string
}

type CreateInput struct {
	UserId  int64
	Title   string
	Content string
	Pics    []string
	MediaId int64
}

type ListFilter struct {
	Status  int // -1=全部(后台); 前台固定1
	Sort    string
	Keyword string
	UserId  int64
	Page    int
	Size    int
}

type IPost interface {
	Create(ctx context.Context, in CreateInput) (int64, error)
	// FrontList 已通过帖子流(sort=new/hot)。
	FrontList(ctx context.Context, f ListFilter) ([]*PostDTO, int, error)
	// Detail 详情(view+1; 非通过状态仅作者可见, viewerId=0 表示未登录)。
	Detail(ctx context.Context, id, viewerId int64) (*PostDTO, error)
	My(ctx context.Context, userId int64, page, size int) ([]*PostDTO, int, error)
	// DeleteOwn 作者软删(status=3)。
	DeleteOwn(ctx context.Context, userId, id int64) error
	// List 后台列表。
	List(ctx context.Context, f ListFilter) ([]*PostDTO, int, error)
	// Audit 审核(仅 0待审 可审, 条件更新幂等)。
	Audit(ctx context.Context, id int64, pass bool, reason string) error
	// Delete 后台硬删(连带评论)。
	Delete(ctx context.Context, id int64) error
}
