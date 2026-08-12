// Package service 推广应用对外接口。
package service

import "context"

type AppDTO struct {
	Id          int64
	Name        string
	Tag         int
	Intro       string
	Avatar      string
	DownloadUrl string
	IosUrl      string
	AndroidUrl  string
	LocIds      []int64
	Rank        int
	DownTotal   int64
	Status      int
	CreatedAt   string
}

type SaveInput struct {
	Id          int64 // 0=新增
	Name        string
	Tag         int
	Intro       string
	Avatar      string
	DownloadUrl string
	IosUrl      string
	AndroidUrl  string
	LocIds      []int64
	Rank        int
	Status      int
}

type ListFilter struct {
	Status  int // -1=全部
	Keyword string
	Page    int
	Size    int
}

type IApplication interface {
	// FrontList 前台: 上架应用(rank desc), loc>0 时按投放位置筛。
	FrontList(ctx context.Context, loc int64) ([]*AppDTO, error)
	// Click 前台: 下载点击计数。
	Click(ctx context.Context, id int64) error
	// List 后台: 分页列表。
	List(ctx context.Context, f ListFilter) ([]*AppDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
}
