// Package service 基础配置对外接口。
package service

import "context"

type ItemDTO struct {
	Id        int64
	Grp       string
	Key       string
	Value     string
	Remark    string
	Status    int
	UpdatedAt string
}

type CreateInput struct {
	Grp    string
	Key    string
	Value  string
	Remark string
	Status int
}

type UpdateInput struct {
	Id     int64
	Grp    string
	Value  string
	Remark string
	Status int
}

type ListFilter struct {
	Grp     string
	Status  int // -1=全部
	Keyword string
	Page    int
	Size    int
}

type IConfig interface {
	// Info 前台: 全量启用配置(key → 原始类型值)。
	Info(ctx context.Context, grp string) (map[string]interface{}, error)
	// List 后台: 分页列表。
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in CreateInput) (int64, error)
	Update(ctx context.Context, in UpdateInput) error
	Delete(ctx context.Context, id int64) error
}
