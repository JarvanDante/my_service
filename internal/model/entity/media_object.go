package entity

import "github.com/gogf/gf/v2/os/gtime"

type MediaObject struct {
	Id          int64       `json:"id"          orm:"id"`
	SiteId      int64       `json:"siteId"      orm:"site_id"`
	Bucket      string      `json:"bucket"      orm:"bucket"`
	ObjectKey   string      `json:"objectKey"   orm:"object_key"`
	Purpose     string      `json:"purpose"     orm:"purpose"`
	ContentType string      `json:"contentType" orm:"content_type"`
	Size        int64       `json:"size"        orm:"size"`
	Etag        string      `json:"etag"        orm:"etag"`
	CreatedBy   int64       `json:"createdBy"   orm:"created_by"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
}
