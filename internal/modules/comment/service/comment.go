// Package service 评论对外接口。
package service

import "context"

type ItemDTO struct {
	Id         int64
	UserId     int64
	ParentId   int64
	RootId     int64
	Content    string
	LikeCount  int
	ReplyCount int
	CreatedAt  string
	Replies    []ItemDTO
}

type AdminItemDTO struct {
	Id         int64
	UserId     int64
	Nickname   string
	Img        string
	IsVip      bool
	MediaType  int
	ContentId  int64
	ParentId   int64
	RootId     int64
	Content    string
	LikeCount  int
	ReplyCount int
	Status     int
	CreatedAt  string
}

type AddInput struct {
	UserId    int64
	MediaType int
	ContentId int64
	ParentId  int64
	Content   string
}

type AdminListFilter struct {
	Status    int    // <0 不过滤
	Kind      string // main / reply / 空
	Keyword   string
	UserId    int64
	MediaType int
	Page      int
	Size      int
}

type IComment interface {
	// Add 发表评论/回复。VIP 直接上墙(status=1); 普通用户待审(status=0)。
	Add(ctx context.Context, in AddInput) (id int64, status int, err error)
	// List 顶层分页 + 每条带全部已上墙回复。
	List(ctx context.Context, mediaType int, contentId int64, page, size int) ([]ItemDTO, int, error)
	// AdminList 后台审核列表(主评+回复)。
	AdminList(ctx context.Context, f AdminListFilter) ([]*AdminItemDTO, int, error)
	// Audit 仅待审可审: pass 上墙并补计数, 拒绝不计数。
	Audit(ctx context.Context, id int64, pass bool) error
}
