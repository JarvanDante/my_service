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

type AddInput struct {
	UserId    int64
	MediaType int
	ContentId int64
	ParentId  int64
	Content   string
}

type IComment interface {
	// Add 发表评论/回复(过敏感词, 回复维护顶层 reply_count, 帖子维护 comment_count)。
	Add(ctx context.Context, in AddInput) (int64, error)
	// List 顶层分页 + 每条带全部回复(root_id 反范式, 两查合并免递归)。
	List(ctx context.Context, mediaType int, contentId int64, page, size int) ([]ItemDTO, int, error)
}
