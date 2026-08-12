// Package service 运营配置对外接口(公告/跳转位/敏感词)。
package service

import "context"

type AnnDTO struct {
	Id        int64
	Title     string
	Content   string
	TextNode  string
	Cover     string
	JumpUrl   string
	SysType   string
	StartAt   string
	EndAt     string
	Status    int
	CreatedAt string
}

type AnnSaveInput struct {
	Id       int64 // 0=新增
	Title    string
	Content  string
	TextNode string
	Cover    string
	JumpUrl  string
	SysType  string
	StartAt  string
	EndAt    string
	Status   int
}

type JtDTO struct {
	Id          int64
	CnName      string
	EnName      string
	Avatar      string
	Link        string
	PicJumpLink string
	Location    int
	Rank        int
	Status      int
	CreatedAt   string
}

type JtSaveInput struct {
	Id          int64 // 0=新增
	CnName      string
	EnName      string
	Avatar      string
	Link        string
	PicJumpLink string
	Location    int
	Rank        int
	Status      int
}

type FwDTO struct {
	Id        int64
	Word      string
	CreatedAt string
}

type PageFilter struct {
	Status   int // -1=全部
	Location int // jumptab 用, 0=全部
	Keyword  string
	Page     int
	Size     int
}

type IOps interface {
	// 前台
	LiveAnnouncements(ctx context.Context, sysType string) ([]AnnDTO, error)
	FrontJumptabs(ctx context.Context, location int) ([]JtDTO, error)
	// 公告管理
	AnnList(ctx context.Context, f PageFilter) ([]*AnnDTO, int, error)
	AnnCreate(ctx context.Context, in AnnSaveInput) (int64, error)
	AnnUpdate(ctx context.Context, in AnnSaveInput) error
	AnnDelete(ctx context.Context, id int64) error
	// 跳转位管理
	JtList(ctx context.Context, f PageFilter) ([]*JtDTO, int, error)
	JtCreate(ctx context.Context, in JtSaveInput) (int64, error)
	JtUpdate(ctx context.Context, in JtSaveInput) error
	JtDelete(ctx context.Context, id int64) error
	// 敏感词管理
	FwList(ctx context.Context, f PageFilter) ([]*FwDTO, int, error)
	FwAdd(ctx context.Context, words []string) (int, error)
	FwDelete(ctx context.Context, id int64) error
}
