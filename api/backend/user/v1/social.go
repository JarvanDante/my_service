// Package v1 后台社交查询接口契约(B6)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type FollowItem struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"` // 关注人
	UserName  string `json:"user_name"`
	HomeId    int64  `json:"home_id"` // 被关注人
	HomeName  string `json:"home_name"`
	CreatedAt string `json:"created_at"`
}

// 关注关系查询
type FollowListReq struct {
	g.Meta `path:"/follows" method:"get" tags:"Backend/Social" summary:"关注关系查询"`
	UserId int64 `json:"user_id" v:"min:0#user_id 不合法"` // 查某人关注了谁
	HomeId int64 `json:"home_id" v:"min:0#home_id 不合法"` // 查谁关注了某人
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type FollowListRes struct {
	List  []FollowItem `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}
