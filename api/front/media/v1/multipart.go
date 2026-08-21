package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type MultipartInitReq struct {
	g.Meta      `path:"/media/multipart/init" method:"post" tags:"Front/Media" summary:"分片上传-初始化(可续传)"`
	Filename    string `json:"filename" v:"required#文件名必填"`
	Purpose     string `json:"purpose" d:"video" dc:"前台仅 video"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size" v:"required|min:1#文件大小必填"`
	PartSize    int64  `json:"part_size"`
}

type MultipartInitRes struct {
	UploadId    string `json:"upload_id"`
	ObjectKey   string `json:"object_key"`
	Bucket      string `json:"bucket"`
	Purpose     string `json:"purpose"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	PartSize    int64  `json:"part_size"`
	PartCount   int    `json:"part_count"`
}

type MultipartPartUploadReq struct {
	g.Meta     `path:"/media/multipart/part" method:"post" tags:"Front/Media" mime:"multipart/form-data" summary:"分片上传-上传一片"`
	UploadId   string            `json:"upload_id" v:"required#upload_id 必填"`
	PartNumber int               `json:"part_number" v:"required|min:1#分片号必填"`
	File       *ghttp.UploadFile `json:"file" type:"file" v:"required#分片文件必填"`
}

type MultipartPartUploadRes struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"etag"`
	Size       int64  `json:"size"`
}

type MultipartPartsReq struct {
	g.Meta   `path:"/media/multipart/parts" method:"get" tags:"Front/Media" summary:"分片上传-已传列表"`
	UploadId string `json:"upload_id" v:"required#upload_id 必填"`
}

type MultipartPartItem struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"etag"`
	Size       int64  `json:"size"`
}

type MultipartPartsRes struct {
	UploadId  string              `json:"upload_id"`
	Status    int                 `json:"status"`
	PartCount int                 `json:"part_count"`
	List      []MultipartPartItem `json:"list"`
}

type MultipartCompleteReq struct {
	g.Meta   `path:"/media/multipart/complete" method:"post" tags:"Front/Media" summary:"分片上传-合并完成"`
	UploadId string                    `json:"upload_id" v:"required#upload_id 必填"`
	Parts    []MultipartCompletePartIn `json:"parts"`
}

type MultipartCompletePartIn struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"etag"`
}

type MultipartCompleteRes struct {
	Id          int64  `json:"id"`
	Url         string `json:"url"`
	ObjectKey   string `json:"object_key"`
	Bucket      string `json:"bucket"`
	Purpose     string `json:"purpose"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type MultipartAbortReq struct {
	g.Meta   `path:"/media/multipart/abort" method:"post" tags:"Front/Media" summary:"分片上传-取消"`
	UploadId string `json:"upload_id" v:"required#upload_id 必填"`
}

type MultipartAbortRes struct{}
