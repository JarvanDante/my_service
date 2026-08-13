// Package service 图集对外接口。
package service

import "context"

type PicDTO struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PhotoDTO struct {
	Id        int64
	Title     string
	Cover     string
	Intro     string
	Category  string
	Tags      []string
	IsVip     int
	Price     float64
	FreeCount int
	Pics      []PicDTO // 前台列表不下发(置 nil), 详情按解锁态截断
	PicCount  int
	ViewCount int64
	BuyCount  int64
	LikeCount int64
	Rank      int
	Status    int
	PublishId int64
	CreatedAt string
	IsBuy     bool
}

// DetailDTO 详情 = 图集 + 解锁判定 + 截断信息。
// PreviewCount/TotalCount 是图集特有的: 未解锁时 Photo.Pics 已被服务端裁到前 FreeCount 张,
// 光看 len(Pics) 前端分不清"这套就这么多图"还是"被付费墙截断了", 所以显式给两个数。
type DetailDTO struct {
	Photo        *PhotoDTO
	Playable     bool
	NeedPay      bool
	NeedVip      bool
	Enough       bool
	Reason       string
	PreviewCount int
	TotalCount   int
}

type SaveInput struct {
	Id        int64
	Title     string
	Cover     string
	Intro     string
	Category  string
	Tags      []string
	IsVip     int
	Price     float64
	FreeCount int
	Pics      []PicDTO
	Rank      int
	Status    int
	PublishId int64
}

type ListFilter struct {
	Category string
	Tag      string
	Keyword  string
	Status   int // -1=全部, 前台固定传 1
	Sort     int // 0综合 1最多观看 2最新 3最多点赞
	Page     int
	Size     int
}

type IPhoto interface {
	// FrontList 只出上架图集; userId>0 时批量标记 is_buy。
	FrontList(ctx context.Context, userId int64, f ListFilter) ([]*PhotoDTO, int, error)
	// Detail 详情并 +1 观看数; 未解锁只返回前 free_count 张预览图。
	Detail(ctx context.Context, userId, id int64) (*DetailDTO, error)
	// Buy 整套购买(事务防重与扣款都在 shared/paywall 里)。
	Buy(ctx context.Context, userId, id int64) (float64, float64, error)

	List(ctx context.Context, f ListFilter) ([]*PhotoDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	Audit(ctx context.Context, id int64, status int) error
}
