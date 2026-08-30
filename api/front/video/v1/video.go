// Package v1 前台视频浏览契约。
// 后台契约在 api/backend/video/v1(路径 /videos), 前台走 /video/xxx, 两边字段裁剪不同:
// 前台不下发 cover_key/source_key/created_by 这类运营侧字段。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CoverUrl    string   `json:"cover_url"`
	SourceUrl   string   `json:"source_url"`
	Category    string   `json:"category"`
	Categories  []string `json:"categories"`
	Duration    int      `json:"duration"`
	CreatedAt   string   `json:"created_at"`
	UpUserId    int64    `json:"up_user_id"`
	UpNickname  string   `json:"up_nickname"`
	UpAvatar    string   `json:"up_avatar"`
	Followed      bool     `json:"followed"`
	CommentCount  int      `json:"comment_count"`
	PreviewSec    int      `json:"preview_sec"`
	NeedVip       bool     `json:"need_vip"`
}

// CategoryListReq 启用中的视频分类, 按权重倒序(公开)。
type CategoryListReq struct {
	g.Meta `path:"/video/categories" method:"get" tags:"Front/Video" summary:"视频分类"`
}
type CategoryListRes struct {
	List []FrontCategoryItem `json:"list"`
}

// ListReq 视频列表(公开, 只出已上架 kind=0)。
// sort: 0综合(人工 sort 权重) 1最新 2时长。
type ListReq struct {
	g.Meta   `path:"/video/list" method:"get" tags:"Front/Video" summary:"视频列表"`
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
	Sort     int    `json:"sort" v:"in:0,1,2#排序方式非法"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// DetailReq 视频详情(公开, 只出已上架)。
type DetailReq struct {
	g.Meta `path:"/video/detail" method:"get" tags:"Front/Video" summary:"视频详情"`
	Id     int64 `json:"id" v:"required|min:1#视频ID必填"`
}
type DetailRes struct {
	Item
}
