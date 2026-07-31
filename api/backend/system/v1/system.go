// Package v1 后台系统模块接口契约(B7): 公告/推送 + 客服配置。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 公告 / 推送 ----------

type NoticeItem struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"` // notice / push
	Status    int    `json:"status"`
	CreatedBy int64  `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

// 发布公告/推送
type PushReq struct {
	g.Meta  `path:"/push" method:"post" tags:"Backend/System" summary:"发布系统公告/推送"`
	Title   string `json:"title"   v:"required#标题必填"`
	Content string `json:"content" v:"required#内容必填"`
	Type    string `json:"type"    v:"required|in:notice,push#类型必填|type 仅支持 notice/push"`
}
type PushRes struct {
	Id int64 `json:"id"`
}

// 公告列表
type NoticeListReq struct {
	g.Meta `path:"/notices" method:"get" tags:"Backend/System" summary:"公告/推送列表"`
	Type   string `json:"type"   v:"in:,notice,push#type 仅支持 notice/push"`
	Status int    `json:"status" v:"in:0,1,2#状态仅支持0/1/2"` // 0全部 1上架 2下线
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type NoticeListRes struct {
	List  []NoticeItem `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// 上/下线公告
type NoticeStatusReq struct {
	g.Meta `path:"/notices/{id}/status" method:"put" tags:"Backend/System" summary:"上/下线公告"`
	Id     int64 `json:"id"     v:"required|min:1#公告ID必填|公告ID必须大于0"`
	Status int   `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type NoticeStatusRes struct{}

// ---------- 客服链接配置 ----------

type CustomerUrlGetReq struct {
	g.Meta `path:"/config/customer-url" method:"get" tags:"Backend/System" summary:"查看客服链接"`
}
type CustomerUrlGetRes struct {
	Url string `json:"url"`
}

type CustomerUrlPutReq struct {
	g.Meta `path:"/config/customer-url" method:"put" tags:"Backend/System" summary:"配置客服链接"`
	Url    string `json:"url" v:"required#链接必填"`
}
type CustomerUrlPutRes struct{}
