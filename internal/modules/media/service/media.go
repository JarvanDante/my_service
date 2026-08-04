package service

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type UploadInput struct {
	File       *ghttp.UploadFile
	Purpose    string
	OperatorId int64
}

type UploadDTO struct {
	Id          int64
	Url         string
	ObjectKey   string
	Bucket      string
	Purpose     string
	ContentType string
	Size        int64
}

type IMedia interface {
	Upload(ctx context.Context, in UploadInput) (*UploadDTO, error)
}
