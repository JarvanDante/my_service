// Package v1 后台系统消息契约(发布/管理)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"` // 0=全员
	Type      int    `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/message" method:"get" tags:"Backend/Message" summary:"消息列表"`
	UserId  string `json:"user_id"` // 空=全部  0=只看全员  >0=指定用户
	Status  string `json:"status"`  // 空=全部  0=撤回  1=发布
	Keyword string `json:"keyword"` // 标题模糊搜索
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta  `path:"/message" method:"post" tags:"Backend/Message" summary:"发布消息(user_id=0 全员)"`
	UserId  int64  `json:"user_id"` // 0=全员, >0=指定用户
	Type    int    `json:"type"`
	Title   string `json:"title" v:"required#标题必填"`
	Content string `json:"content" v:"required#内容必填"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta  `path:"/message/{id}" method:"put" tags:"Backend/Message" summary:"更新/撤回消息"`
	Id      int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int    `json:"status" v:"in:0,1#状态非法"` // 0=撤回 1=发布
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/message/{id}" method:"delete" tags:"Backend/Message" summary:"删除消息"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
