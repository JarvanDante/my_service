package service

import "context"

type VideoDTO struct {
	Id            int64
	Title         string
	Description   string
	CoverUrl      string
	CoverKey      string
	CoverMediaId  int64
	SourceUrl     string
	SourceKey     string
	SourceMediaId int64
	MediaCode     string
	Category      string
	Categories    []string
	Tags          []string
	Duration      int
	Sort          int
	Status        int
	UpUserId      int64
	UpNickname    string
	CreatedBy     int64
	CreatedAt     string
	UpdatedAt     string
}

type ListInput struct {
	Keyword   string
	MediaCode string
	Kind      int
	Status    int
	Page      int
	Size      int
}

type ListDTO struct {
	List  []*VideoDTO
	Total int
	Page  int
	Size  int
}

// FrontListInput 前台列表入参。刻意不带 Status —— 前台永远只看已上架,
// 不给调用方留下"传个 status 就能看到草稿"的口子。
type FrontListInput struct {
	Keyword  string
	Category string
	Tag      string
	Tags     []string
	Kind     int // 0视频 2动漫
	Sort     int // 0综合(sort权重) 1最新 2时长
	Page     int
	Size     int
}

type SaveInput struct {
	Id            int64
	Title         string
	Description   string
	CoverUrl      string
	CoverKey      string
	CoverMediaId  int64
	SourceUrl     string
	SourceKey     string
	SourceMediaId int64
	MediaCode     string
	Kind          int
	Category      string
	Categories    []string
	Tags          []string
	Duration      int
	Sort          int
	Status        int
	UpUserId      int64
	OperatorId    int64
}

type IVideo interface {
	// FrontList 前台列表: 只出 status=1 已上架。
	FrontList(ctx context.Context, in FrontListInput) (*ListDTO, error)
	// FrontDetail 前台详情: 未上架一律当作不存在(不泄露"存在但下架了")。
	FrontDetail(ctx context.Context, id int64) (*VideoDTO, error)

	List(ctx context.Context, in ListInput) (*ListDTO, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status int) error
	ListMediaAssets(ctx context.Context, page, size int, keyword string, kind int) ([]MediaAssetDTO, int, error)
	PickMedia(ctx context.Context, code string, operatorId int64, kind int) (int64, error)
	SyncMedia(ctx context.Context, operatorId int64, kind int) (*SyncMediaDTO, error)
}

type MediaAssetDTO struct {
	Id          string
	Title       string
	CoverUrl    string
	PlayUrl     string
	DurationSec int
	Picked      bool
	LocalId     int64
}

type SyncMediaDTO struct {
	Created int
	Updated int
	Total   int
}
