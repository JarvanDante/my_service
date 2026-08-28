// Package service 小说对外接口。
package service

import "context"

type NovelDTO struct {
	Id           int64
	Title        string
	Author       string
	Cover        string
	Intro        string
	Category     string
	Tags         []string
	IsVip        int
	Price        float64
	FreeChapter  int
	ChapterCount int
	WordCount    int64
	IsAudio      int
	ViewCount    int64
	BuyCount     int64
	LikeCount    int64
	UpdateStatus int
	Rank         int
	Status       int
	PublishId    int64
	CreatedAt    string
	IsBuy        bool
}

type ChapterDTO struct {
	Id        int64
	NovelId   int64
	Seq       int
	Title     string
	Content   string
	WordCount int
	AudioUrl  string
	Status    int
	CreatedAt string
	IsFree    bool
	Playable  bool
}

// DetailDTO 详情 = 作品 + 解锁判定。
type DetailDTO struct {
	Novel    *NovelDTO
	Playable bool
	NeedPay  bool
	NeedVip  bool
	Enough   bool
	Reason   string
}

type ReadDTO struct {
	ChapterId int64
	NovelId   int64
	Seq       int
	Title     string
	Content   string
	WordCount int
	AudioUrl  string
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
	Tags         []string
	IsVip        int
	Price        float64
	FreeChapter  int
	IsAudio      int
	UpdateStatus int
	Rank         int
	Status       int
	PublishId    int64
}

type ChapterInput struct {
	Id       int64
	NovelId  int64
	Seq      int
	Title    string
	Content  string
	AudioUrl string
	Status   int
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

type INovel interface {
	// FrontList 只出上架作品; userId>0 时批量标记 is_buy。
	FrontList(ctx context.Context, userId int64, f ListFilter) ([]*NovelDTO, int, error)
	// FrontCategories 已上架作品里出现过的分类名, 去重保序。
	FrontCategories(ctx context.Context) ([]string, error)
	// Detail 详情并 +1 观看数。
	Detail(ctx context.Context, userId, id int64) (*DetailDTO, error)
	// Chapters 目录, 逐章标出 is_free / playable。
	Chapters(ctx context.Context, userId, id int64, desc bool) (string, []*ChapterDTO, error)
	// Read 读章节, 未解锁直接报错(免费章节不需要登录)。
	Read(ctx context.Context, userId, chapterId int64) (*ReadDTO, error)
	// Buy 整部购买(事务: 购买记录唯一约束防重 + 条件扣款 + 流水 + 销量+1)。
	Buy(ctx context.Context, userId, id int64) (float64, float64, error)
	MayLike(ctx context.Context, id int64, size int) ([]*NovelDTO, error)

	List(ctx context.Context, f ListFilter) ([]*NovelDTO, int, error)
	Create(ctx context.Context, in SaveInput) (int64, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	Audit(ctx context.Context, id int64, status int) error

	ChapterList(ctx context.Context, novelId int64, page, size int) ([]*ChapterDTO, int, error)
	ChapterCreate(ctx context.Context, in ChapterInput) (int64, error)
	ChapterUpdate(ctx context.Context, in ChapterInput) error
	ChapterDelete(ctx context.Context, id int64) error
}
