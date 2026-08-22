// Package v1 前台帖子契约(移植自 tianbi post/community 核心面)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id           int64    `json:"id"`
	UserId       int64    `json:"user_id"`
	Nickname     string   `json:"nickname"`
	Img          string   `json:"img"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Pics         []string `json:"pics"`
	Topics       []string `json:"topics"`
	VideoUrl     string   `json:"video_url"`
	MediaId      int64    `json:"media_id"`
	ViewCount    int64    `json:"view_count"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	// Status 审核态: 0待审 1通过 2拒绝 3已删。
	// 这里**不能**用 omitempty —— 待审正是 0, omitempty 会把它整个字段抹掉,
	// 前端拿到 undefined 就没法区分"待审"和"没这个字段", 是踩过的坑。
	// 公开流只出已通过的帖子, 多下发一个 status=1 无害, 换来"我的帖子"能正确显示审核态。
	Status       int    `json:"status"`
	RejectReason string `json:"reject_reason"`
	CreatedAt    string `json:"created_at"`
}

// CreateReq 发帖(需登录, 过敏感词, 待审核)。
type CreateReq struct {
	g.Meta   `path:"/post/create" method:"post" tags:"Front/Post" summary:"发帖"`
	Title    string   `json:"title" v:"required|max-length:128#标题必填|标题过长"`
	Content  string   `json:"content" v:"required|max-length:5000#内容必填|内容过长"`
	Pics     []string `json:"pics" v:"required#请上传图片"`
	Topics   []string `json:"topics" v:"required#请选择话题"`
	VideoUrl string   `json:"video_url"`
	MediaId  int64    `json:"media_id"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// ListReq 帖子流(公开, 仅已通过; sort=new/hot)。
type ListReq struct {
	g.Meta  `path:"/post/list" method:"get" tags:"Front/Post" summary:"帖子列表"`
	Sort     string `json:"sort"`     // new(默认)/hot
	Keyword  string `json:"keyword"`  // 标题模糊
	UserId   int64  `json:"user_id"`  // 0=全部, >0 该用户已通过帖子
	Category string `json:"category"` // 普通分类名(匹配 topics)
	Follow   int    `json:"follow"`   // 1=只看当前用户关注的人
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// DetailReq 详情(公开, 浏览数+1; 待审/拒绝仅作者可见)。
type DetailReq struct {
	g.Meta `path:"/post/detail" method:"get" tags:"Front/Post" summary:"帖子详情"`
	Id     int64 `json:"id" v:"required|min:1#帖子ID必填"`
}
type DetailRes struct {
	Post Item `json:"post"`
}

// MyReq 我的帖子(需登录, 含待审/拒绝)。
type MyReq struct {
	g.Meta `path:"/post/my" method:"get" tags:"Front/Post" summary:"我的帖子"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type MyRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

// DeleteReq 删除自己的帖子(软删)。
type DeleteReq struct {
	g.Meta `path:"/post/delete" method:"post" tags:"Front/Post" summary:"删除我的帖子"`
	Id     int64 `json:"id" v:"required|min:1#帖子ID必填"`
}
type DeleteRes struct{}
