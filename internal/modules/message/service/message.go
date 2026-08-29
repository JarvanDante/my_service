// Package service 系统消息对外接口。
package service

import "context"

type MsgDTO struct {
	Id        int64
	UserId    int64
	Type      int
	Title     string
	Content   string
	IsRead    bool
	Status    int
	CreatedAt string
}

type CreateInput struct {
	UserId  int64 // 0=全员
	Type    int
	Title   string
	Content string
}

type UpdateInput struct {
	Id      int64
	Type    int
	Title   string
	Content string
	Status  int
}

type ListFilter struct {
	UserId  int64 // 后台: -1=全部, 0=只看全员, >0=指定用户
	Status  int   // -1=全部
	Keyword string
	Page    int
	Size    int
}

type InteractDTO struct {
	Id            int64
	Channel       string
	SubType       string
	IsRead        bool
	CreatedAt     string
	ActorId       int64
	ActorName     string
	ActorAvatar   string
	ActorSex      int
	ActorCount    int
	MediaType     int
	ContentId     int64
	ObjectTitle   string
	TargetType    string
	CommentId     int64
	RootCommentId int64
	Snippet       string
}

type UnreadBreakdown struct {
	Sys     int
	Comment int
	Like    int
}

type IMessage interface {
	// MyList 前台: 全员消息 + 发给我的(status=1), 带已读标记。
	MyList(ctx context.Context, userId int64, page, size int) ([]MsgDTO, int, error)
	// UnreadCount 前台: 站内信未读数(不含互动)。
	UnreadCount(ctx context.Context, userId int64) (int, error)
	// UnreadAll 前台: 站内信 + 评论/点赞未读。
	UnreadAll(ctx context.Context, userId int64) (UnreadBreakdown, error)
	// MarkRead 前台: 标记已读(id>0 单条; all=true 全部)。
	MarkRead(ctx context.Context, userId, id int64, all bool) error
	InteractList(ctx context.Context, userId int64, channel string, page, size int) ([]InteractDTO, int, error)
	MarkInteractRead(ctx context.Context, userId, id int64, all bool, channel string) error
	// List 后台: 分页列表。
	List(ctx context.Context, f ListFilter) ([]*MsgDTO, int, error)
	Create(ctx context.Context, in CreateInput) (int64, error)
	Update(ctx context.Context, in UpdateInput) error
	Delete(ctx context.Context, id int64) error
}
