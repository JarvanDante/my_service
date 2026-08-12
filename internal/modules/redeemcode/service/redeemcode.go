// Package service 兑换码对外接口。
package service

import "context"

// UseResult 兑换成功结果。
type UseResult struct {
	Desc string // 如: 兑换金币x100
}

// MyRecordDTO 前台我的兑换记录。
type MyRecordDTO struct {
	Code      string
	Desc      string
	ActivedAt string
}

// ItemDTO 后台兑换码。
type ItemDTO struct {
	Id         int64
	Name       string
	Code       string
	CardType   int
	Value      int
	TotalTimes int
	UsedTimes  int
	Status     int
	ExpiredAt  string
	CreatedAt  string
}

// RecordDTO 后台使用记录。
type RecordDTO struct {
	Id        int64
	UserId    int64
	CodeId    int64
	Code      string
	Name      string
	CardType  int
	Value     int
	CreatedAt string
}

type CreateInput struct {
	Name       string
	Code       string // 空则自动生成
	Value      int
	TotalTimes int
	ExpiredAt  string
	Status     int
}

type UpdateInput struct {
	Id         int64
	Name       string
	Value      int
	TotalTimes int
	ExpiredAt  string
	Status     int
}

type ListFilter struct {
	Status  int // -1=全部  0=禁用  1=启用
	Keyword string
	Page    int
	Size    int
}

type RecordFilter struct {
	UserId int64
	Code   string
	Page   int
	Size   int
}

type IRedeemCode interface {
	// Use 前台: 使用兑换码(事务发金币, 防重用/防超发)。
	Use(ctx context.Context, userId int64, code string) (*UseResult, error)
	// MyRecords 前台: 我的兑换记录。
	MyRecords(ctx context.Context, userId int64, page, size int) ([]MyRecordDTO, error)
	// List 后台: 兑换码分页列表。
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in CreateInput) (int64, string, error)
	Update(ctx context.Context, in UpdateInput) error
	Delete(ctx context.Context, id int64) error
	// Records 后台: 使用记录分页列表。
	Records(ctx context.Context, f RecordFilter) ([]*RecordDTO, int, error)
}
