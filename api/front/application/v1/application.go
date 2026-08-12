// Package v1 前台推广应用契约(移植自 tianbi ApplicationInfo)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type AppItem struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Tag         int     `json:"tag"`
	Desc        string  `json:"desc"` // 对客户端保持 tianbi 字段名
	Avatar      string  `json:"avatar"`
	DownloadUrl string  `json:"download_url"`
	IosUrl      string  `json:"iosUrl"`
	AndroidUrl  string  `json:"androidUrl"`
	LocIds      []int64 `json:"locIds"`
	Rank        int     `json:"rank"`
	DownTotal   int64   `json:"downTotal"`
}

// ListReq 推广应用列表(公开)。
type ListReq struct {
	g.Meta `path:"/application/list" method:"get" tags:"Front/Application" summary:"推广应用列表"`
	Loc    int64 `json:"loc"` // 可选: 按投放位置筛
}
type ListRes struct {
	List []AppItem `json:"list"`
}

// ClickReq 下载点击上报(公开, 计数)。
type ClickReq struct {
	g.Meta `path:"/application/click" method:"post" tags:"Front/Application" summary:"下载点击上报"`
	Id     int64 `json:"id" v:"required|min:1#应用ID必填"`
}
type ClickRes struct{}
