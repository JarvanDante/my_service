// Package service 短剧对外接口。
package service

import "context"

type DramaDTO struct {
	Id           int64
	Title        string
	Cover        string
	Intro        string
	Category     string
	Categories   []string
	Tags         []string
	IsVip        int
	Price        float64
	FreeEpisode  int
	EpisodeCount int
	ViewCount    int64
	BuyCount     int64
	LikeCount    int64
	UpdateStatus int
	Rank         int
	IsRecommend  int
	Status       int
	MediaCode    string
	CreatedAt    string
}

// EpisodeDTO 剧集；前台字段按当前用户解锁态填充。
type EpisodeDTO struct {
	Id         int64
	DramaId    int64
	Seq        int
	Title      string
	MediaCode  string
	Duration   int
	Cover      string
	PayMode    int
	Price      float64
	Status     int
	MediaReady int
	CreatedAt  string
	IsFree     bool
	IsVip      bool
	IsBuy      bool
	Playable   bool
	NeedPay    bool
	NeedVip    bool
}

type PlayDTO struct {
	Episode    *EpisodeDTO
	PlayUrl    string
	PreviewSec int
	Enough     bool
	Reason     string
	PrevId     int64
	NextId     int64
}

type SaveInput struct {
	Id           int64
	Title        string
	Cover        string
	Intro        string
	Category     string
	Categories   []string
	Tags         []string
	IsVip        int
	Price        float64
	FreeEpisode  int
	UpdateStatus int
	Rank         int
	IsRecommend  int
	Status       int
}

type EpisodeInput struct {
	Id      int64
	PayMode int
	Price   float64
	Status  int
}

type ListFilter struct {
	Category      string
	Categories    []string
	Tag           string
	Tags          []string
	Keyword       string
	Status        int // -1=全部, 前台固定传 1
	Sort          int // 0综合 1最多观看 2最新 3最多点赞
	PayType       int // 0全部 1VIP 2金币 3免费(按整部默认)
	UpdateStatus  int // 0不限 1连载 2完结
	Ids           []int64
	ExcludeIds    []int64
	Shuffle       bool
	OnlyRecommend bool
	Page          int
	Size          int
}

type SyncDTO struct {
	Created int
	Updated int
	Missing int
	Total   int
}

type MediaDramaDTO struct {
	Id           string
	Title        string
	CoverUrl     string
	Intro        string
	ChapterCount int
	EpisodeCount int
	UpdateStatus int
	Picked       bool
	LocalId      int64
}

type IDrama interface {
	// FrontList 只出上架短剧。
	FrontList(ctx context.Context, f ListFilter) ([]*DramaDTO, int, error)
	// Detail 上架短剧详情并 +1 观看数。
	Detail(ctx context.Context, id int64) (*DramaDTO, error)
	// Episodes 可播剧集, 逐集标出当前用户的解锁态。
	Episodes(ctx context.Context, userId, id int64, desc bool) (string, []*EpisodeDTO, error)
	// Play 签发播放地址, 未解锁只给试看地址。
	Play(ctx context.Context, userId, episodeId int64) (*PlayDTO, error)
	// Buy 单集购买(唯一约束防重 + 条件扣款 + 流水)。
	Buy(ctx context.Context, userId, episodeId int64) (float64, float64, error)

	List(ctx context.Context, f ListFilter) ([]*DramaDTO, int, error)
	Update(ctx context.Context, in SaveInput) error
	Delete(ctx context.Context, id int64) error
	Audit(ctx context.Context, id int64, status int) error
	EpisodeList(ctx context.Context, dramaId int64, page, size int) ([]*EpisodeDTO, int, error)
	EpisodeUpdate(ctx context.Context, in EpisodeInput) error
	SyncEpisodes(ctx context.Context, id int64) (*SyncDTO, error)

	ListMediaDramas(ctx context.Context, page, size int, keyword string) ([]MediaDramaDTO, int, error)
	PickMedia(ctx context.Context, code string) (int64, error)
}
