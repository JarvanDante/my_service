// Package v1 后台媒体上传接口契约。
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadReq 上传媒体到 MinIO。
type UploadReq struct {
	g.Meta  `path:"/media/upload" method:"post" tags:"Backend/Media" mime:"multipart/form-data" summary:"上传媒体文件"`
	File    *ghttp.UploadFile `json:"file" type:"file" v:"required#文件必填"`
	Purpose string            `json:"purpose" d:"image" dc:"用途: image/cover/video/avatar"`
}

type UploadRes struct {
	Id          int64  `json:"id"`
	Url         string `json:"url"`
	ObjectKey   string `json:"object_key"`
	Bucket      string `json:"bucket"`
	Purpose     string `json:"purpose"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}
