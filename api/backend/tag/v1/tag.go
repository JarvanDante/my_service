// Package v1 后台标签契约(CRUD 维护)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id          int64  `json:"id"`
	ContentType int    `json:"content_type"`
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type ListReq struct {
	g.Meta      `path:"/tag" method:"get" tags:"Backend/Tag" summary:"标签列表"`
	ContentType int    `json:"content_type"` // 0/空=全部
	Status      string `json:"status"`       // 空=全部  0=只看禁用  1=只看启用
	Keyword     string `json:"keyword"`      // 标签名模糊
	Page        int    `json:"page"`
	Size        int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta      `path:"/tag" method:"post" tags:"Backend/Tag" summary:"新增标签"`
	ContentType int    `json:"content_type" v:"required|in:1,2,3,4,5,6,7#内容类型必填|内容类型非法"`
	Name        string `json:"name" v:"required#标签名必填"`
	Rank        int    `json:"rank"`
	Status      int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/tag/{id}" method:"put" tags:"Backend/Tag" summary:"更新标签"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#标签ID必填"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/tag/{id}" method:"delete" tags:"Backend/Tag" summary:"删除标签"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#标签ID必填"`
}
type DeleteRes struct{}
