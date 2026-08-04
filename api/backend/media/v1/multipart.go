package v1

import "github.com/gogf/gf/v2/frame/g"

// MultipartInitReq 初始化分片上传。
type MultipartInitReq struct {
	g.Meta      `path:"/media/multipart/init" method:"post" tags:"Backend/Media" summary:"分片上传-初始化"`
	Filename    string `json:"filename" v:"required#文件名必填"`
	Purpose     string `json:"purpose" d:"video" dc:"用途: video/cover/image/avatar"`
	ContentType string `json:"content_type" dc:"MIME, 可空则按扩展名推断"`
	Size        int64  `json:"size" v:"required|min:1#文件大小必填"`
	PartSize    int64  `json:"part_size" dc:"分片字节数, 默认 8MiB, 最小 5MiB(末片除外)"`
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

// MultipartPresignReq 批量获取分片预签名 URL。
type MultipartPresignReq struct {
	g.Meta      `path:"/media/multipart/presign" method:"post" tags:"Backend/Media" summary:"分片上传-预签名"`
	UploadId    string `json:"upload_id" v:"required#upload_id 必填"`
	PartNumbers []int  `json:"part_numbers" v:"required#part_numbers 必填" dc:"分片号, 从 1 起"`
}

type MultipartPresignItem struct {
	PartNumber int    `json:"part_number"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	ExpiresIn  int64  `json:"expires_in"`
}

type MultipartPresignRes struct {
	List []MultipartPresignItem `json:"list"`
}

// MultipartPartsReq 查询已上传分片(断点续传)。
type MultipartPartsReq struct {
	g.Meta   `path:"/media/multipart/parts" method:"get" tags:"Backend/Media" summary:"分片上传-已传列表"`
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

// MultipartCompleteReq 合并分片。parts 可空: 空则服务端 ListParts 后合并。
type MultipartCompleteReq struct {
	g.Meta   `path:"/media/multipart/complete" method:"post" tags:"Backend/Media" summary:"分片上传-合并完成"`
	UploadId string                    `json:"upload_id" v:"required#upload_id 必填"`
	Parts    []MultipartCompletePartIn `json:"parts" dc:"可选; 空则用 MinIO 已传列表"`
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

// MultipartAbortReq 取消分片上传。
type MultipartAbortReq struct {
	g.Meta   `path:"/media/multipart/abort" method:"post" tags:"Backend/Media" summary:"分片上传-取消"`
	UploadId string `json:"upload_id" v:"required#upload_id 必填"`
}

type MultipartAbortRes struct{}
