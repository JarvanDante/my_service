package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UploadReq struct {
	g.Meta  `path:"/media/upload" method:"post" tags:"Front/Media" mime:"multipart/form-data" summary:"上传图片或视频"`
	File    *ghttp.UploadFile `json:"file" type:"file" v:"required#文件必填"`
	Purpose string            `json:"purpose" d:"image" dc:"用途: image/avatar/video/ad"`
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

type ObjectReq struct {
	g.Meta    `path:"/media/object" method:"get" tags:"Front/Media" summary:"读取本站已上传图片(密文原样)"`
	Url       string `json:"u" in:"query" dc:"MinIO 地址"`
	ObjectKey string `json:"object_key" in:"query"`
}

type ObjectRes struct{}
