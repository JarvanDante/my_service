// Package v1 前台运营配置契约(公告 + 跳转位, 公开)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type AnnItem struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	TextNode string `json:"textNode"`
	Cover    string `json:"cover"`
	JumpUrl  string `json:"jumpUrl"`
	SysType  string `json:"sysType"`
	StartAt  string `json:"start"`
	EndAt    string `json:"end"`
}

// AnnouncementReq 有效期内的启用公告(公开)。
type AnnouncementReq struct {
	g.Meta  `path:"/announcement/list" method:"post" tags:"Front/Ops" summary:"官方公告"`
	SysType string `json:"sys_type"` // 空=全部, app/h5app
}
type AnnouncementRes struct {
	List []AnnItem `json:"list"`
}

type JumptabItem struct {
	Id          int64  `json:"id"`
	CnName      string `json:"cnName"`
	EnName      string `json:"enName"`
	Avatar      string `json:"avatar"`
	Link        string `json:"link"`
	PicJumpLink string `json:"picJumpLink"`
	Location    int    `json:"location"`
	Rank        int    `json:"rank"`
}

// JumptabReq 跳转位列表(公开)。
type JumptabReq struct {
	g.Meta   `path:"/jumptab/list" method:"get" tags:"Front/Ops" summary:"跳转位列表"`
	Location int `json:"location"` // 0=全部
}
type JumptabRes struct {
	List []JumptabItem `json:"list"`
}
