// Package service 系统模块对外接口(B7)。
package service

import "context"

type NoticeDTO struct {
	Id        int64
	Title     string
	Content   string
	Type      string
	Status    int
	CreatedBy int64
	CreatedAt string
}

type PushInput struct {
	Title      string
	Content    string
	Type       string // notice / push
	OperatorId int64
}

type NoticeListInput struct {
	Type   string
	Status int // 0全部 1上架 2下线
	Page   int
	Size   int
}

type NoticeListDTO struct {
	List  []*NoticeDTO
	Total int
	Page  int
	Size  int
}

type ISystem interface {
	// 公告 / 推送
	Push(ctx context.Context, in PushInput) (int64, error)
	Notices(ctx context.Context, in NoticeListInput) (*NoticeListDTO, error)
	SetNoticeStatus(ctx context.Context, id int64, status int) error
	// 配置
	GetCustomerUrl(ctx context.Context) (string, error)
	SetCustomerUrl(ctx context.Context, url string) error
}
