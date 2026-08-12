// Package v1 后台意见反馈契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id          int64    `json:"id"`
	UserId      int64    `json:"user_id"`
	Type        int      `json:"type"`
	ProblemType int      `json:"problem_type"`
	Content     string   `json:"content"`
	Pics        []string `json:"pics"`
	SysInfo     string   `json:"sys_info"`
	MediaId     int64    `json:"media_id"`
	MediaTitle  string   `json:"media_title"`
	Status      int      `json:"status"`
	Reply       string   `json:"reply"`
	CreatedAt   string   `json:"created_at"`
}

type ListReq struct {
	g.Meta `path:"/feedback" method:"get" tags:"Backend/Feedback" summary:"反馈列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
	Status int `json:"status"` // 0全部 1处理中 2已处理
	Type   int `json:"type"`   // 0全部
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type HandleReq struct {
	g.Meta `path:"/feedback/{id}/handle" method:"post" tags:"Backend/Feedback" summary:"处理反馈(回复+置为已处理)"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#反馈ID必填"`
	Reply  string `json:"reply"`
	Status int    `json:"status" v:"in:1,2#状态非法"`
}
type HandleRes struct{}
