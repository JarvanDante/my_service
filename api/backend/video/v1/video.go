// Package v1 后台视频管理接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type VideoItem struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CoverUrl      string `json:"cover_url"`
	CoverKey      string `json:"cover_key"`
	CoverMediaId  int64  `json:"cover_media_id"`
	SourceUrl     string `json:"source_url"`
	SourceKey     string `json:"source_key"`
	SourceMediaId int64  `json:"source_media_id"`
	Duration      int    `json:"duration"`
	Sort          int    `json:"sort"`
	Status        int    `json:"status"`
	CreatedBy     int64  `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ListReq struct {
	g.Meta  `path:"/videos" method:"get" tags:"Backend/Video" summary:"视频列表"`
	Keyword string `json:"keyword"`
	Status  int    `json:"status" d:"9" v:"in:0,1,2,9#状态不合法" dc:"9=全部 0草稿 1上架 2下架"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []VideoItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type CreateReq struct {
	g.Meta        `path:"/videos" method:"post" tags:"Backend/Video" summary:"新增视频"`
	Title         string `json:"title" v:"required#标题必填"`
	Description   string `json:"description"`
	CoverUrl      string `json:"cover_url"`
	CoverKey      string `json:"cover_key"`
	CoverMediaId  int64  `json:"cover_media_id"`
	SourceUrl     string `json:"source_url" v:"required#请先上传视频"`
	SourceKey     string `json:"source_key"`
	SourceMediaId int64  `json:"source_media_id"`
	Duration      int    `json:"duration" v:"min:0#时长不合法"`
	Sort          int    `json:"sort"`
	Status        int    `json:"status" v:"in:0,1,2#status 仅支持 0/1/2"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta        `path:"/videos/{id}" method:"put" tags:"Backend/Video" summary:"编辑视频"`
	Id            int64  `json:"id" v:"required|min:1#ID必填"`
	Title         string `json:"title" v:"required#标题必填"`
	Description   string `json:"description"`
	CoverUrl      string `json:"cover_url"`
	CoverKey      string `json:"cover_key"`
	CoverMediaId  int64  `json:"cover_media_id"`
	SourceUrl     string `json:"source_url" v:"required#请先上传视频"`
	SourceKey     string `json:"source_key"`
	SourceMediaId int64  `json:"source_media_id"`
	Duration      int    `json:"duration" v:"min:0#时长不合法"`
	Sort          int    `json:"sort"`
	Status        int    `json:"status" v:"in:0,1,2#status 仅支持 0/1/2"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/videos/{id}" method:"delete" tags:"Backend/Video" summary:"删除视频"`
	Id     int64 `json:"id" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

type StatusReq struct {
	g.Meta `path:"/videos/{id}/status" method:"put" tags:"Backend/Video" summary:"视频上下架"`
	Id     int64 `json:"id" v:"required|min:1#ID必填"`
	Status int   `json:"status" v:"in:0,1,2#status 仅支持 0/1/2"`
}
type StatusRes struct{}
