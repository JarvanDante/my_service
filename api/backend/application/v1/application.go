// Package v1 后台推广应用契约(CRUD)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Tag         int     `json:"tag"`
	Intro       string  `json:"intro"`
	Avatar      string  `json:"avatar"`
	DownloadUrl string  `json:"download_url"`
	IosUrl      string  `json:"ios_url"`
	AndroidUrl  string  `json:"android_url"`
	LocIds      []int64 `json:"loc_ids"`
	Rank        int     `json:"rank"`
	DownTotal   int64   `json:"down_total"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/application" method:"get" tags:"Backend/Application" summary:"推广应用列表"`
	Status  string `json:"status"`  // 空=全部  0=下架  1=上架
	Keyword string `json:"keyword"` // 名称模糊搜索
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta      `path:"/application" method:"post" tags:"Backend/Application" summary:"新增推广应用"`
	Name        string  `json:"name" v:"required#应用名必填"`
	Tag         int     `json:"tag"`
	Intro       string  `json:"intro"`
	Avatar      string  `json:"avatar"`
	DownloadUrl string  `json:"download_url"`
	IosUrl      string  `json:"ios_url"`
	AndroidUrl  string  `json:"android_url"`
	LocIds      []int64 `json:"loc_ids"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta      `path:"/application/{id}" method:"put" tags:"Backend/Application" summary:"更新推广应用"`
	Id          int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name        string  `json:"name"`
	Tag         int     `json:"tag"`
	Intro       string  `json:"intro"`
	Avatar      string  `json:"avatar"`
	DownloadUrl string  `json:"download_url"`
	IosUrl      string  `json:"ios_url"`
	AndroidUrl  string  `json:"android_url"`
	LocIds      []int64 `json:"loc_ids"`
	Rank        int     `json:"rank"`
	Status      int     `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/application/{id}" method:"delete" tags:"Backend/Application" summary:"删除推广应用"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
