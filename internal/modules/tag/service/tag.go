// Package service 标签对外接口。
package service

import "context"

// RepoItem 前台标签(id + name)。
type RepoItem struct {
	Id   int64
	Name string
}

type ItemDTO struct {
	Id          int64
	ContentType int
	Name        string
	Rank        int
	Status      int
	CreatedAt   string
}

type CreateInput struct {
	ContentType int
	Name        string
	Rank        int
	Status      int
}

type UpdateInput struct {
	Id     int64
	Name   string
	Rank   int
	Status int
}

type ListFilter struct {
	ContentType int
	Status      int // -1=全部  0=只看禁用  1=只看启用
	Keyword     string
	Page        int
	Size        int
}

type ITag interface {
	// Repo 前台: 按内容类型取启用标签(rank desc)。
	Repo(ctx context.Context, contentType, page, size int) ([]RepoItem, error)
	// List 后台: 分页列表。
	List(ctx context.Context, f ListFilter) ([]*ItemDTO, int, error)
	Create(ctx context.Context, in CreateInput) (int64, error)
	Update(ctx context.Context, in UpdateInput) error
	Delete(ctx context.Context, id int64) error
}
