// Package v1 前台小说契约(与前台漫画契约同构, 差异在章节载体是正文/音频而非图片)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64    `json:"id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip"`
	Price        float64  `json:"price"`
	FreeChapter  int      `json:"free_chapter"`
	ChapterCount int      `json:"chapter_count"`
	WordCount    int64    `json:"word_count"`
	IsAudio      int      `json:"is_audio"`
	ViewCount    int64    `json:"view_count"`
	LikeCount    int64    `json:"like_count"`
	UpdateStatus int      `json:"update_status"`
	IsBuy        bool     `json:"is_buy"`
	CreatedAt    string   `json:"created_at"`
}

// ListReq 小说列表(公开)。sort: 0综合(rank) 1最多观看 2最新 3最多点赞。
type ListReq struct {
	g.Meta   `path:"/novel/list" method:"get" tags:"Front/Novel" summary:"小说列表"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
	Keyword  string `json:"keyword"`
	Sort     int    `json:"sort"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// DetailReq 小说详情(公开; 带登录态则返回解锁信息)。
type DetailReq struct {
	g.Meta `path:"/novel/detail" method:"get" tags:"Front/Novel" summary:"小说详情"`
	Id     int64 `json:"id" v:"required|min:1#小说ID必填"`
}
type DetailRes struct {
	Item
	Playable bool   `json:"playable"` // 整部是否已解锁
	NeedPay  bool   `json:"need_pay"`
	NeedVip  bool   `json:"need_vip"`
	Enough   bool   `json:"enough"` // 余额是否够买
	Reason   string `json:"reason"`
}

type ChapterItem struct {
	Id        int64  `json:"id"`
	Seq       int    `json:"seq"`
	Title     string `json:"title"`
	WordCount int    `json:"word_count"`
	HasAudio  bool   `json:"has_audio"` // 目录不下发音频地址, 只标是否有
	IsFree    bool   `json:"is_free"`   // 是否属于免费章节
	Playable  bool   `json:"playable"`  // 当前用户是否可读
}

// ChaptersReq 章节目录(公开)。
type ChaptersReq struct {
	g.Meta `path:"/novel/chapters" method:"get" tags:"Front/Novel" summary:"小说章节目录"`
	Id     int64 `json:"id" v:"required|min:1#小说ID必填"`
	Desc   bool  `json:"desc"` // true=倒序
}
type ChaptersRes struct {
	NovelId int64         `json:"novel_id"`
	Title   string        `json:"title"`
	List    []ChapterItem `json:"list"`
}

// ReadReq 读某一章(免费章节公开; 付费章节需已购/VIP)。
type ReadReq struct {
	g.Meta    `path:"/novel/read" method:"get" tags:"Front/Novel" summary:"阅读小说章节"`
	ChapterId int64 `json:"chapter_id" v:"required|min:1#章节ID必填"`
}
type ReadRes struct {
	ChapterId int64  `json:"chapter_id"`
	NovelId   int64  `json:"novel_id"`
	Seq       int    `json:"seq"`
	Title     string `json:"title"`
	Content   string `json:"content"` // 正文
	WordCount int    `json:"word_count"`
	AudioUrl  string `json:"audio_url"` // 有声小说才有值
	PrevId    int64  `json:"prev_id"`
	NextId    int64  `json:"next_id"`
}

// BuyReq 整部购买(需登录)。价格服务端定, 不接受客户端传金额。
type BuyReq struct {
	g.Meta `path:"/novel/buy" method:"post" tags:"Front/Novel" summary:"购买小说"`
	Id     int64 `json:"id" v:"required|min:1#小说ID必填"`
}
type BuyRes struct {
	Price   float64 `json:"price"`
	Balance float64 `json:"balance"` // 购买后余额
}

// MayLikeReq 猜你喜欢(公开, 同分类随机取)。
type MayLikeReq struct {
	g.Meta `path:"/novel/maylike" method:"get" tags:"Front/Novel" summary:"猜你喜欢"`
	Id     int64 `json:"id"`
	Size   int   `json:"size"`
}
type MayLikeRes struct {
	List []Item `json:"list"`
}
