// Package service 推广/兑换码模块对外接口(B3)。
package service

import "context"

type CodeDTO struct {
	Id        int64
	Name      string
	Code      string
	CodeKey   string
	Type      string
	ObjectId  int64
	AddNum    int
	CanUseNum int
	UsedNum   int
	Status    int // 0未使用 1已使用 -1作废
	ExpiredAt int64
	CreatedAt string
}

type CodeListInput struct {
	Keyword string
	CodeKey string
	Type    string
	Status  int
	Page    int
	Size    int
}

// CodeGenInput 批量生成入参。
type CodeGenInput struct {
	Name      string
	Type      string // point 金币 / group 用户组
	ObjectId  int64  // type=group 时为用户组ID
	AddNum    int    // 金币数 / 天数
	CanUseNum int    // 每码可用次数, 默认1
	Count     int    // 生成数量 1~1000
	ExpiredAt int64  // 过期 epoch 秒, 0 不过期
}

type CodeGenDTO struct {
	CodeKey string // 批次号
	Count   int
	Codes   []string // 生成的码
}

type CodeLogDTO struct {
	Id        int64
	CodeId    int64
	Code      string
	Name      string
	Type      string
	UserId    int64
	Username  string
	AddNum    int
	CreatedAt string
}

type CodeLogListInput struct {
	CodeId int64
	UserId int64
	Code   string
	Page   int
	Size   int
}

type PageDTO[T any] struct {
	List  []*T
	Total int
	Page  int
	Size  int
}

type IPromo interface {
	Codes(ctx context.Context, in CodeListInput) (*PageDTO[CodeDTO], error)
	GenerateCodes(ctx context.Context, in CodeGenInput) (*CodeGenDTO, error)
	VoidCode(ctx context.Context, id int64) error
	CodeLogs(ctx context.Context, in CodeLogListInput) (*PageDTO[CodeLogDTO], error)
}
