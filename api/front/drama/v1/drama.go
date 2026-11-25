// Package v1 前台短剧契约(列表 / 详情 / 选集 / 播放 / 单集购买)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64    `json:"id"`
	Title        string   `json:"title"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip"`
	Price        float64  `json:"price"`
	FreeEpisode  int      `json:"free_episode"`
	EpisodeCount int      `json:"episode_count"`
	ViewCount    int64    `json:"view_count"`
	LikeCount    int64    `json:"like_count"`
	UpdateStatus int      `json:"update_status"`
	CreatedAt    string   `json:"created_at"`
}

type CategoryItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Kind int    `json:"kind"`
}

// CategoryListReq 启用中的短剧分类, 按权重倒序(公开)。
type CategoryListReq struct {
	g.Meta `path:"/drama/categories" method:"get" tags:"Front/Drama" summary:"短剧分类"`
}
type CategoryListRes struct {
	List []CategoryItem `json:"list"`
}

// ListReq 短剧列表(公开)。sort: 0综合 1最多观看 2最新 3最多点赞; pay_type 按整部默认: 0全部 1VIP 2金币 3免费。
type ListReq struct {
	g.Meta     `path:"/drama/list" method:"get" tags:"Front/Drama" summary:"短剧列表"`
	Category   string `json:"category"`
	Categories string `json:"categories"`
	Tag        string `json:"tag"`
	Tags       string `json:"tags"`
	Keyword    string `json:"keyword"`
	Sort       int    `json:"sort"`
	PayType    int    `json:"pay_type"`
	Recommend  int    `json:"recommend"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type DetailReq struct {
	g.Meta `path:"/drama/detail" method:"get" tags:"Front/Drama" summary:"短剧详情"`
	Id     int64 `json:"id" v:"required|min:1#短剧ID必填"`
}
type DetailRes struct {
	Item
}

type EpisodeItem struct {
	Id       int64   `json:"id"`
	Seq      int     `json:"seq"`
	Title    string  `json:"title"`
	Duration int     `json:"duration"`
	Cover    string  `json:"cover"`
	PayMode  int     `json:"pay_mode"`
	IsFree   bool    `json:"is_free"`
	IsVip    bool    `json:"is_vip"`
	Price    float64 `json:"price"`
	IsBuy    bool    `json:"is_buy"`
	Playable bool    `json:"playable"`
	NeedPay  bool    `json:"need_pay"`
	NeedVip  bool    `json:"need_vip"`
}

// EpisodesReq 选集列表(公开; 带登录态逐集返回解锁态)。
type EpisodesReq struct {
	g.Meta `path:"/drama/episodes" method:"get" tags:"Front/Drama" summary:"短剧选集"`
	Id     int64 `json:"id" v:"required|min:1#短剧ID必填"`
	Desc   bool  `json:"desc"`
}
type EpisodesRes struct {
	DramaId int64         `json:"drama_id"`
	Title   string        `json:"title"`
	List    []EpisodeItem `json:"list"`
}

// PlayReq 播放某一集；未解锁时下发 5 秒试看地址。
type PlayReq struct {
	g.Meta    `path:"/drama/play" method:"get" tags:"Front/Drama" summary:"播放短剧剧集"`
	EpisodeId int64 `json:"episode_id" v:"required|min:1#剧集ID必填"`
}
type PlayRes struct {
	EpisodeId  int64   `json:"episode_id"`
	DramaId    int64   `json:"drama_id"`
	Seq        int     `json:"seq"`
	Title      string  `json:"title"`
	Cover      string  `json:"cover"`
	Duration   int     `json:"duration"`
	PlayUrl    string  `json:"play_url"`
	PreviewSec int     `json:"preview_sec"`
	Playable   bool    `json:"playable"`
	IsBuy      bool    `json:"is_buy"`
	NeedPay    bool    `json:"need_pay"`
	NeedVip    bool    `json:"need_vip"`
	Price      float64 `json:"price"`
	Enough     bool    `json:"enough"`
	Reason     string  `json:"reason"`
	PrevId     int64   `json:"prev_id"`
	NextId     int64   `json:"next_id"`
}

// BuyReq 单集购买(需登录)，价格服务端定。
type BuyReq struct {
	g.Meta    `path:"/drama/buy" method:"post" tags:"Front/Drama" summary:"购买短剧单集"`
	EpisodeId int64 `json:"episode_id" v:"required|min:1#剧集ID必填"`
}
type BuyRes struct {
	Price   float64 `json:"price"`
	Balance float64 `json:"balance"`
}
