// Package v1 后台评论审核契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Nickname   string `json:"nickname"`
	Img        string `json:"img"`
	IsVip      bool   `json:"is_vip"`
	MediaType  int    `json:"media_type"`
	ContentId  int64  `json:"content_id"`
	ParentId   int64  `json:"parent_id"`
	RootId     int64  `json:"root_id"`
	Content    string `json:"content"`
	LikeCount  int    `json:"like_count"`
	ReplyCount int    `json:"reply_count"`
	Status      int    `json:"status"` // 0待审 1已上墙 2已拒绝
	BelongLabel string `json:"belong_label"`
	CreatedAt   string `json:"created_at"`
}

type ListReq struct {
	g.Meta    `path:"/comment" method:"get" tags:"Backend/Comment" summary:"评论审核列表"`
	Status    string `json:"status"`     // 空=全部 0待审 1已上墙 2已拒绝
	Kind      string `json:"kind"`       // 空=全部 main=主评 reply=回复
	Keyword   string `json:"keyword"`    // 内容模糊
	UserId    int64  `json:"user_id"`    // 0=全部
	MediaType int    `json:"media_type"` // 0=全部 1视频 2帖子 4漫画 7小说 8动漫 9抖音（8/9 仅筛选，仍存 media_type=1）
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type AuditReq struct {
	g.Meta `path:"/comment/{id}/audit" method:"post" tags:"Backend/Comment" summary:"审核评论"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
	Pass   bool  `json:"pass"`
}
type AuditRes struct{}
