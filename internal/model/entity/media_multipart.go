package entity

import "github.com/gogf/gf/v2/os/gtime"

const (
	MultipartStatusUploading = 0
	MultipartStatusCompleted = 1
	MultipartStatusAborted   = 2
)

type MediaMultipart struct {
	Id            int64       `json:"id"            orm:"id"`
	SiteId        int64       `json:"siteId"        orm:"site_id"`
	UploadId      string      `json:"uploadId"      orm:"upload_id"`
	MinioUploadId string      `json:"minioUploadId" orm:"minio_upload_id"`
	Bucket        string      `json:"bucket"        orm:"bucket"`
	ObjectKey     string      `json:"objectKey"     orm:"object_key"`
	Purpose       string      `json:"purpose"       orm:"purpose"`
	Filename      string      `json:"filename"      orm:"filename"`
	ContentType   string      `json:"contentType"   orm:"content_type"`
	Size          int64       `json:"size"          orm:"size"`
	PartSize      int64       `json:"partSize"      orm:"part_size"`
	PartCount     int         `json:"partCount"     orm:"part_count"`
	Status        int         `json:"status"        orm:"status"`
	MediaId       int64       `json:"mediaId"       orm:"media_id"`
	CreatedBy     int64       `json:"createdBy"     orm:"created_by"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"`
}
