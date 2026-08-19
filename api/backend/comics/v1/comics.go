// Package v1 后台漫画契约(作品 CRUD + 审核 + 章节 CRUD)。
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
	ViewCount    int64    `json:"view_count"`
	BuyCount     int64    `json:"buy_count"`
	LikeCount    int64    `json:"like_count"`
	UpdateStatus int      `json:"update_status"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status"`
	PublishId    int64    `json:"publish_id"`
	MediaCode    string   `json:"media_code"`
	CreatedAt    string   `json:"created_at"`
}

type ListReq struct {
	g.Meta   `path:"/comics" method:"get" tags:"Backend/Comics" summary:"漫画列表"`
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

type CreateReq struct {
	g.Meta       `path:"/comics" method:"post" tags:"Backend/Comics" summary:"新增漫画"`
	Title        string   `json:"title" v:"required#标题必填"`
	Author       string   `json:"author"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price        float64  `json:"price"`
	FreeChapter  int      `json:"free_chapter"`
	UpdateStatus int      `json:"update_status" v:"in:0,1,2#连载状态非法"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status" v:"in:0,1,2#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta       `path:"/comics/{id}" method:"put" tags:"Backend/Comics" summary:"编辑漫画"`
	Id           int64    `json:"id" in:"path" v:"required|min:1#ID必填"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price        float64  `json:"price"`
	FreeChapter  int      `json:"free_chapter"`
	UpdateStatus int      `json:"update_status" v:"in:0,1,2#连载状态非法"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status" v:"in:0,1,2#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/comics/{id}" method:"delete" tags:"Backend/Comics" summary:"删除漫画(连带章节)"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// AuditReq 上下架/审核: status 0待上架 1上架 2下架。
type AuditReq struct {
	g.Meta `path:"/comics/{id}/audit" method:"post" tags:"Backend/Comics" summary:"漫画上下架"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
	Status int   `json:"status" v:"in:0,1,2#状态非法"`
}
type AuditRes struct{}

type ChapterItem struct {
	Id        int64  `json:"id"`
	ComicsId  int64  `json:"comics_id"`
	Seq       int    `json:"seq"`
	Title     string `json:"title"`
	Pics      []Pic  `json:"pics"`
	PicCount  int    `json:"pic_count"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type Pic struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ChaptersReq struct {
	g.Meta `path:"/comics/{id}/chapters" method:"get" tags:"Backend/Comics" summary:"章节列表"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#漫画ID必填"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type ChaptersRes struct {
	List  []ChapterItem `json:"list"`
	Total int           `json:"total"`
}

type ChapterCreateReq struct {
	g.Meta `path:"/comics/{id}/chapters" method:"post" tags:"Backend/Comics" summary:"新增章节"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#漫画ID必填"`
	Seq    int    `json:"seq" v:"required|min:1#章节序号必填"`
	Title  string `json:"title"`
	Pics   []Pic  `json:"pics"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type ChapterCreateRes struct {
	Id int64 `json:"id"`
}

type ChapterUpdateReq struct {
	g.Meta `path:"/comics-chapters/{id}" method:"put" tags:"Backend/Comics" summary:"编辑章节"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Seq    int    `json:"seq"`
	Title  string `json:"title"`
	Pics   []Pic  `json:"pics"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type ChapterUpdateRes struct{}

type ChapterDeleteReq struct {
	g.Meta `path:"/comics-chapters/{id}" method:"delete" tags:"Backend/Comics" summary:"删除章节"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type ChapterDeleteRes struct{}

type MediaComicsItem struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	CoverUrl     string `json:"cover_url"`
	Intro        string `json:"intro"`
	ChapterCount int    `json:"chapter_count"`
	Picked       bool   `json:"picked"`
	LocalId      int64  `json:"local_id"`
}

type MediaComicsListReq struct {
	g.Meta  `path:"/media-comics" method:"get" tags:"Backend/Comics" summary:"媒资中心可选用漫画"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type MediaComicsListRes struct {
	List  []MediaComicsItem `json:"list"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

type MediaComicsPickReq struct {
	g.Meta `path:"/media-comics/{id}/pick" method:"post" tags:"Backend/Comics" summary:"选用漫画媒资并写入本站列表"`
	Id     string `json:"id" v:"required#媒资ID必填"`
}
type MediaComicsPickRes struct {
	Id int64 `json:"id"`
}
