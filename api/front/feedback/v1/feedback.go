// Package v1 前台意见反馈契约(移植自 tianbi feedback)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type AddReq struct {
	g.Meta      `path:"/feedback/add" method:"post" tags:"Front/Feedback" summary:"提交意见反馈"`
	Type        int      `json:"type"`         // 1用户反馈 2程序反馈
	ProblemType int      `json:"problem_type"` // 问题类型
	Content     string   `json:"content" v:"required#反馈内容必填"`
	Pics        []string `json:"pics"`
	SysInfo     string   `json:"sys_info"`
	MediaId     int64    `json:"media_id"`
	MediaTitle  string   `json:"media_title"`
}
type AddRes struct {
	Id int64 `json:"id"`
}
