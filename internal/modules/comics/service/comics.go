// Package service 漫画对外接口。
package service

import "context"

type PicDTO struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Key    string `json:"key,omitempty"`
}

type ComicsDTO struct {
	Id           int64
	Title        string
	Author       string
	Cover        string
	Intro        string
	Category     string
	Categories   []string
	Tags         []string
	IsVip        int
	Price        float64
	FreeChapter  int
	ChapterCount int
	ViewCount    int64
	BuyCount     int64
	LikeCount    int64
	UpdateStatus int
	Rank         int
	IsRecommend  int
	Status       int
	PublishId    int64
	MediaCode    string
	CreatedAt    string
	IsBuy        bool
}

type ChapterDTO struct {
	Id        int64
	ComicsId  int64
	Seq       int
	Title     string
	Pics      []PicDTO
	PicCount  int
	Status    int
	CreatedAt string
	IsFree    bool
	Playable  bool
}

// DetailDTO 详情 = 作品 + 解锁判定。
type DetailDTO struct {
	Comics   *ComicsDTO
	Playable bool
	NeedPay  bool
	NeedVip  bool
	Enough   bool
	Reason   string
}

type ReadDTO struct {
	ChapterId int64
	ComicsId  int64
	Seq       int
	Title     string
	Pics      []PicDTO
	PrevId    int64
	NextId    int64
}

type SaveInput struct {
	Id           int64
	Title        string
	Author       string
	Cover        string
	Intro        string
	Category     string
	Categories   []string
	Tags         []string
	IsVip        int
	Price        float64
	FreeChapter  int
	UpdateStatus int
	Rank         int
	IsRecommend  int
	Status       int
	PublishId    int64
}

type ChapterInput struct {
	Id       int64
	ComicsId int64
	Seq      int
	Title    string
	Pics     []PicDTO
	Status   int
}

type ListFilter struct {
	Category      string
	Categories    []string
	Tag           string
	Tags          []string
	Keyword       string
	Status        int  // -1=全部, 前台固定传 1
	Sort          int  // 0综合(推荐权重) 1最多观看 2最新 3最多点赞
	PayType       int  // 0全部 1VIP 2付费解锁 3免费
	OnlyRecommend bool // 仅 is_recommend=1, 给 H5 推荐栏
	Page          int
	Size          int
}

type IComics interface {
	// FrontList 只出上架作品; userId>0 时批量标记 is_buy。
	FrontList(ctx context.Context, userId int64, f ListFilter) ([]*ComicsDTO, int, error)
	// Detail 详情并 +1 观看数。
	Detail(ctx context.Context, userId, id int64) (*DetailDTO, error)
	// Chapters 目录, 逐章标出 is_free / playable。
	Chapters(ctx context.Context, userId, id int64, desc bool) (string, []*ChapterDTO, error)
	// Read 读章节, 未解锁直接报错(免费章节不需要登录)。
	Read(ctx context.Context, userId, chapterId int64) (*ReadDTO, error)
	// Buy 整部购买(事务: 购买记录唯一约束防重 + 条件扣款 + 流水 + 销量+1)。
	Buy(ctx context.Context, userId, id int64) (float64, float64, error)
	MayLike(ctx context.Context, id int64, size int) ([]*ComicsDTO, error)

	List(ctx context.Context, f ListFilter) ([]*ComicsDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	Audit(ctx context.Context, id int64, status int) error

	ChapterList(ctx context.Context, comicsId int64, page, size int) ([]*ChapterDTO, int, error)
	ChapterCreate(ctx context.Context, in ChapterInput) (int64, error)
	ChapterUpdate(ctx context.Context, in ChapterInput) error
	ChapterDelete(ctx context.Context, id int64) error

	ListMediaComics(ctx context.Context, page, size int, keyword string) ([]MediaComicsDTO, int, error)
	PickMedia(ctx context.Context, code string) (int64, error)
}

type MediaComicsDTO struct {
	Id           string
	Title        string
	CoverUrl     string
	Intro        string
	ChapterCount int
	Picked       bool
	LocalId      int64
}
