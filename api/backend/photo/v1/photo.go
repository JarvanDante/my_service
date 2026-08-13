// Package v1 后台图集契约(CRUD + 上下架)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Pic struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Item struct {
	Id        int64    `json:"id"`
	Title     string   `json:"title"`
	Cover     string   `json:"cover"`
	Intro     string   `json:"intro"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	IsVip     int      `json:"is_vip"`
	Price     float64  `json:"price"`
	FreeCount int      `json:"free_count"`
	Pics      []Pic    `json:"pics"` // 后台要能核对全部图片, 不截断
	PicCount  int      `json:"pic_count"`
	ViewCount int64    `json:"view_count"`
	BuyCount  int64    `json:"buy_count"`
	LikeCount int64    `json:"like_count"`
	Rank      int      `json:"rank"`
	Status    int      `json:"status"`
	PublishId int64    `json:"publish_id"`
	CreatedAt string   `json:"created_at"`
}

type ListReq struct {
	g.Meta   `path:"/photos" method:"get" tags:"Backend/Photo" summary:"图集列表"`
	Status   string `json:"status"` // 空=全部
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta    `path:"/photos" method:"post" tags:"Backend/Photo" summary:"新增图集"`
	Title     string   `json:"title" v:"required#标题必填"`
	Cover     string   `json:"cover"`
	Intro     string   `json:"intro"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	IsVip     int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price     float64  `json:"price"`
	FreeCount int      `json:"free_count"`
	Pics      []Pic    `json:"pics"`
	Rank      int      `json:"rank"`
	Status    int      `json:"status" v:"in:0,1,2#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta    `path:"/photos/{id}" method:"put" tags:"Backend/Photo" summary:"编辑图集"`
	Id        int64    `json:"id" in:"path" v:"required|min:1#ID必填"`
	Title     string   `json:"title"`
	Cover     string   `json:"cover"`
	Intro     string   `json:"intro"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	IsVip     int      `json:"is_vip" v:"in:0,1#VIP标记非法"`
	Price     float64  `json:"price"`
	FreeCount int      `json:"free_count"`
	Pics      []Pic    `json:"pics"` // 传了才整体覆盖(同时刷新 pic_count)
	Rank      int      `json:"rank"`
	Status    int      `json:"status" v:"in:0,1,2#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/photos/{id}" method:"delete" tags:"Backend/Photo" summary:"删除图集"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// AuditReq 上下架/审核: status 0待上架 1上架 2下架。
type AuditReq struct {
	g.Meta `path:"/photos/{id}/audit" method:"post" tags:"Backend/Photo" summary:"图集上下架"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
	Status int   `json:"status" v:"in:0,1,2#状态非法"`
}
type AuditRes struct{}
