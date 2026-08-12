// Package service 意见反馈对外接口。
package service

import "context"

type AddInput struct {
	UserId      int64
	Type        int
	ProblemType int
	Content     string
	Pics        []string
	SysInfo     string
	MediaId     int64
	MediaTitle  string
}

type ItemDTO struct {
	Id          int64
	UserId      int64
	Type        int
	ProblemType int
	Content     string
	Pics        []string
	SysInfo     string
	MediaId     int64
	MediaTitle  string
	Status      int
	Reply       string
	CreatedAt   string
}

type ListFilter struct {
	Page   int
	Size   int
	Status int
	Type   int
}

type IFeedback interface {
	Add(ctx context.Context, in AddInput) (int64, error)
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Handle(ctx context.Context, id int64, reply string, status int) error
}
