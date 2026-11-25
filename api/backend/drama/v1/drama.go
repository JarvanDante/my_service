// Package v1 后台短剧契约(短剧 / 剧集 / 媒资选用 / 同步剧集)。
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
	BuyCount     int64    `json:"buy_count"`
	LikeCount    int64    `json:"like_count"`
	UpdateStatus int      `json:"update_status"`
	Rank         int      `json:"rank"`
	IsRecommend  int      `json:"is_recommend"`
	Status       int      `json:"status"`
	MediaCode    string   `json:"media_code"`
	CreatedAt    string   `json:"created_at"`
}

type ListReq struct {
	g.Meta   `path:"/dramas" method:"get" tags:"Backend/Drama" summary:"短剧列表"`
	Status   string `json:"status"` // 空=全部
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type UpdateReq struct {
	g.Meta       `path:"/dramas/{id}" method:"put" tags:"Backend/Drama" summary:"编辑短剧"`
	Id           int64    `json:"id" in:"path" v:"required|min:1#ID必填"`
	Title        string   `json:"title"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price        float64  `json:"price"`
	FreeEpisode  int      `json:"free_episode"`
	UpdateStatus int      `json:"update_status" v:"in:0,1,2#连载状态非法"`
	Rank         int      `json:"rank"`
	IsRecommend  int      `json:"is_recommend" v:"in:0,1#推荐标记非法"`
	Status       int      `json:"status" v:"in:0,1,2#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/dramas/{id}" method:"delete" tags:"Backend/Drama" summary:"删除短剧(连带本站剧集)"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// AuditReq 上下架: status 0待上架 1上架 2下架。
type AuditReq struct {
	g.Meta `path:"/dramas/{id}/audit" method:"post" tags:"Backend/Drama" summary:"短剧上下架"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
	Status int   `json:"status" v:"in:0,1,2#状态非法"`
}
type AuditRes struct{}

type EpisodeItem struct {
	Id         int64   `json:"id"`
	DramaId    int64   `json:"drama_id"`
	Seq        int     `json:"seq"`
	Title      string  `json:"title"`
	MediaCode  string  `json:"media_code"`
	Duration   int     `json:"duration"`
	Cover      string  `json:"cover"`
	PayMode    int     `json:"pay_mode"`
	Price      float64 `json:"price"`
	Status     int     `json:"status"`
	MediaReady int     `json:"media_ready"`
	CreatedAt  string  `json:"created_at"`
}

type EpisodesReq struct {
	g.Meta `path:"/dramas/{id}/episodes" method:"get" tags:"Backend/Drama" summary:"剧集列表"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#短剧ID必填"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type EpisodesRes struct {
	List  []EpisodeItem `json:"list"`
	Total int           `json:"total"`
}

// EpisodeUpdateReq 单集付费与上下架: pay_mode 0跟随整部 1免费 2金币 3VIP。
type EpisodeUpdateReq struct {
	g.Meta  `path:"/drama-episodes/{id}" method:"put" tags:"Backend/Drama" summary:"编辑剧集"`
	Id      int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	PayMode int     `json:"pay_mode" v:"in:0,1,2,3#付费方式非法"`
	Price   float64 `json:"price"`
	Status  int     `json:"status" v:"in:0,1#状态非法"`
}
type EpisodeUpdateRes struct{}

type SyncReq struct {
	g.Meta `path:"/dramas/{id}/sync-episodes" method:"post" tags:"Backend/Drama" summary:"按集号同步剧集"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#短剧ID必填"`
}
type SyncRes struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Missing int `json:"missing"`
	Total   int `json:"total"`
}

type MediaDramaItem struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	CoverUrl     string `json:"cover_url"`
	Intro        string `json:"intro"`
	ChapterCount int    `json:"chapter_count"`
	EpisodeCount int    `json:"episode_count"`
	UpdateStatus int    `json:"update_status"`
	Picked       bool   `json:"picked"`
	LocalId      int64  `json:"local_id"`
}

type MediaDramaListReq struct {
	g.Meta  `path:"/media-dramas" method:"get" tags:"Backend/Drama" summary:"媒资中心可选用短剧"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type MediaDramaListRes struct {
	List  []MediaDramaItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

type MediaDramaPickReq struct {
	g.Meta `path:"/media-dramas/{id}/pick" method:"post" tags:"Backend/Drama" summary:"选用短剧媒资并同步剧集"`
	Id     string `json:"id" v:"required#媒资ID必填"`
}
type MediaDramaPickRes struct {
	Id int64 `json:"id"`
}
