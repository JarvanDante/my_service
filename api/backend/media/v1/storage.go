package v1

import "github.com/gogf/gf/v2/frame/g"

type StorageInitReq struct {
	g.Meta      `path:"/media/storage/init" method:"post" tags:"Backend/Media" summary:"站点封面等-统一存储预签名"`
	Filename    string `json:"filename" v:"required#文件名必填"`
	Purpose     string `json:"purpose" d:"cover" dc:"cover/image/avatar/ad/video"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size" v:"required|min:1#文件大小必填"`
}
type StorageInitRes struct {
	Id          string `json:"id"`
	UploadUrl   string `json:"upload_url"`
	Method      string `json:"method"`
	Bucket      string `json:"bucket"`
	ObjectKey   string `json:"object_key"`
	ExpireSec   int    `json:"expire_sec"`
	PublicUrl   string `json:"public_url"`
	ContentType string `json:"content_type"`
}

type StorageConfirmReq struct {
	g.Meta `path:"/media/storage/confirm" method:"post" tags:"Backend/Media" summary:"站点封面等-确认统一存储上传"`
	Id     string `json:"id" v:"required#对象ID必填"`
}
type StorageConfirmRes struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	ObjectKey string `json:"object_key"`
	Bucket    string `json:"bucket"`
	Size      int64  `json:"size"`
}
