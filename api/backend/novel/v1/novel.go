// Package v1 后台小说契约(作品 CRUD + 审核 + 章节 CRUD)。
// 路径用复数 /novels, 与前台 /novel/... 拉开语义, 避免看日志时混淆前后台。
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
	BuyCount     int64    `json:"buy_count"`
	LikeCount    int64    `json:"like_count"`
	UpdateStatus int      `json:"update_status"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status"`
	PublishId    int64    `json:"publish_id"`
	CreatedAt    string   `json:"created_at"`
}

type ListReq struct {
	g.Meta   `path:"/novels" method:"get" tags:"Backend/Novel" summary:"小说列表"`
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
	g.Meta       `path:"/novels" method:"post" tags:"Backend/Novel" summary:"新增小说"`
	Title        string   `json:"title" v:"required#标题必填"`
	Author       string   `json:"author"`
	Cover        string   `json:"cover"`
	Intro        string   `json:"intro"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	IsVip        int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price        float64  `json:"price"`
	FreeChapter  int      `json:"free_chapter"`
	IsAudio      int      `json:"is_audio" v:"in:0,1#有声标记非法"`
	UpdateStatus int      `json:"update_status" v:"in:0,1,2#连载状态非法"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status" v:"in:0,1,2#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta       `path:"/novels/{id}" method:"put" tags:"Backend/Novel" summary:"编辑小说"`
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
	IsAudio      int      `json:"is_audio" v:"in:0,1#有声标记非法"`
	UpdateStatus int      `json:"update_status" v:"in:0,1,2#连载状态非法"`
	Rank         int      `json:"rank"`
	Status       int      `json:"status" v:"in:0,1,2#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/novels/{id}" method:"delete" tags:"Backend/Novel" summary:"删除小说(连带章节)"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// AuditReq 上下架/审核: status 0待上架 1上架 2下架。
type AuditReq struct {
	g.Meta `path:"/novels/{id}/audit" method:"post" tags:"Backend/Novel" summary:"小说上下架"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
	Status int   `json:"status" v:"in:0,1,2#状态非法"`
}
type AuditRes struct{}

type ChapterItem struct {
	Id        int64  `json:"id"`
	NovelId   int64  `json:"novel_id"`
	Seq       int    `json:"seq"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
	AudioUrl  string `json:"audio_url"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ChaptersReq struct {
	g.Meta `path:"/novels/{id}/chapters" method:"get" tags:"Backend/Novel" summary:"章节列表"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#小说ID必填"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type ChaptersRes struct {
	List  []ChapterItem `json:"list"`
	Total int           `json:"total"`
}

// ChapterCreateReq 不收 word_count: 字数由正文实时算, 免得人工填错后统计失真。
type ChapterCreateReq struct {
	g.Meta   `path:"/novels/{id}/chapters" method:"post" tags:"Backend/Novel" summary:"新增章节"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#小说ID必填"`
	Seq      int    `json:"seq" v:"required|min:1#章节序号必填"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AudioUrl string `json:"audio_url"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type ChapterCreateRes struct {
	Id int64 `json:"id"`
}

type ChapterUpdateReq struct {
	g.Meta   `path:"/novel-chapters/{id}" method:"put" tags:"Backend/Novel" summary:"编辑章节"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Seq      int    `json:"seq"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AudioUrl string `json:"audio_url"`
	Status   int    `json:"status" v:"in:0,1#状态非法"`
}
type ChapterUpdateRes struct{}

type ChapterDeleteReq struct {
	g.Meta `path:"/novel-chapters/{id}" method:"delete" tags:"Backend/Novel" summary:"删除章节"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type ChapterDeleteRes struct{}
