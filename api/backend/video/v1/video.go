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
	MediaCode     string   `json:"media_code"`
	Category      string   `json:"category"`
	Categories    []string `json:"categories"`
	Tags          []string `json:"tags"`
	Duration      int      `json:"duration"`
	Sort          int      `json:"sort"`
	Status        int      `json:"status"`
	SubmitSource  int      `json:"submit_source"`
	RejectReason  string   `json:"reject_reason"`
	UpUserId      int64    `json:"up_user_id"`
	UpNickname    string   `json:"up_nickname"`
	CreatedBy     int64    `json:"created_by"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type ListReq struct {
	g.Meta  `path:"/videos" method:"get" tags:"Backend/Video" summary:"视频列表"`
	Keyword   string `json:"keyword"`
	MediaCode string `json:"media_code" dc:"媒资短码, 精确匹配"`
	Kind         int `json:"kind" d:"0" dc:"0视频 2动漫 3抖音"`
	Status       int `json:"status" d:"9" v:"in:0,1,2,3,4,8,9#状态不合法" dc:"9=全部 0草稿 1上架 2下架 3待审 4拒绝 8=草稿+下架"`
	SubmitSource int `json:"submit_source" d:"9" v:"in:0,1,9#来源不合法" dc:"9=全部 0后台 1用户"`
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
	SourceUrl     string `json:"source_url"`
	SourceKey     string `json:"source_key"`
	SourceMediaId int64  `json:"source_media_id"`
	MediaCode     string   `json:"media_code"`
	Kind          int      `json:"kind" d:"0" dc:"0视频 2动漫 3抖音"`
	Category      string   `json:"category"`
	Categories    []string `json:"categories"`
	Tags          []string `json:"tags"`
	Duration      int      `json:"duration" v:"min:0#时长不合法"`
	Sort          int      `json:"sort"`
	Status        int      `json:"status" v:"in:0,1,2#status 仅支持 0/1/2"`
	UpUserId      int64    `json:"up_user_id"`
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
	SourceUrl     string `json:"source_url"`
	SourceKey     string `json:"source_key"`
	SourceMediaId int64  `json:"source_media_id"`
	MediaCode     string   `json:"media_code"`
	Category      string   `json:"category"`
	Categories    []string `json:"categories"`
	Tags          []string `json:"tags"`
	Duration      int      `json:"duration" v:"min:0#时长不合法"`
	Sort          int      `json:"sort"`
	Status        int      `json:"status" v:"in:0,1,2,3,4#status 不合法"`
	UpUserId      int64    `json:"up_user_id"`
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

type AuditReq struct {
	g.Meta       `path:"/videos/{id}/audit" method:"post" tags:"Backend/Video" summary:"审核用户上传抖音"`
	Id           int64  `json:"id" v:"required|min:1#ID必填"`
	Pass         bool   `json:"pass"`
	RejectReason string `json:"reject_reason"`
}
type AuditRes struct{}

type MediaAssetItem struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	CoverUrl    string `json:"cover_url"`
	PlayUrl     string `json:"play_url"`
	DurationSec int    `json:"duration_sec"`
	Picked      bool   `json:"picked"`
	LocalId     int64  `json:"local_id"`
}

type MediaAssetListReq struct {
	g.Meta  `path:"/media-assets" method:"get" tags:"Backend/Video" summary:"媒资中心可选用列表"`
	Keyword string `json:"keyword"`
	Kind    int    `json:"kind" d:"0" dc:"0视频 2动漫 3抖音"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type MediaAssetListRes struct {
	List  []MediaAssetItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

type MediaPickReq struct {
	g.Meta `path:"/media-assets/{id}/pick" method:"post" tags:"Backend/Video" summary:"选用媒资并写入视频列表"`
	Id     string `json:"id" v:"required#媒资ID必填"`
	Kind   int    `json:"kind" d:"0" dc:"0视频 2动漫 3抖音"`
}
type MediaPickRes struct {
	Id int64 `json:"id"`
}

type SyncMediaReq struct {
	g.Meta `path:"/videos/sync-media" method:"post" tags:"Backend/Video" summary:"从媒资中心同步视频"`
	Kind   int `json:"kind" d:"0" dc:"0视频 2动漫 3抖音"`
}
type SyncMediaRes struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Total   int `json:"total"`
}
