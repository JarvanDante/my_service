// Package service 收藏/点赞对外接口。
package service

import "context"

type ItemDTO struct {
	ContentId int64
	MediaType int
	CreatedAt string
}

type OperateInput struct {
	UserId    int64
	Ids       []int64
	MediaType int
	Flag      bool // true=添加 false=取消
	Type      int  // 1收藏 2点赞 3踩
}

type ListFilter struct {
	UserId    int64
	Type      int
	MediaType int // 0=全部
	Page      int
	Size      int
}

type ICollect interface {
	// Operate 添加/取消 收藏/点赞/踩(幂等: 重复添加不报错)。
	Operate(ctx context.Context, in OperateInput) error
	// Delete 批量取消。
	Delete(ctx context.Context, userId int64, ids []int64, mediaType, opType int) error
	// List 我的列表。
	List(ctx context.Context, f ListFilter) ([]ItemDTO, int, error)
}
