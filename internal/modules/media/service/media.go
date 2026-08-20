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

type MultipartInitInput struct {
	Filename    string
	Purpose     string
	ContentType string
	Size        int64
	PartSize    int64
	OperatorId  int64
}

type MultipartInitDTO struct {
	UploadId    string
	ObjectKey   string
	Bucket      string
	Purpose     string
	ContentType string
	Size        int64
	PartSize    int64
	PartCount   int
}

type MultipartPresignInput struct {
	UploadId    string
	PartNumbers []int
	OperatorId  int64
}

type MultipartPresignItemDTO struct {
	PartNumber int
	Url        string
	Method     string
	ExpiresIn  int64
}

type MultipartPartDTO struct {
	PartNumber int
	Etag       string
	Size       int64
}

type MultipartPartsDTO struct {
	UploadId  string
	Status    int
	PartCount int
	List      []MultipartPartDTO
}

type MultipartCompletePartIn struct {
	PartNumber int
	Etag       string
}

type MultipartCompleteInput struct {
	UploadId   string
	Parts      []MultipartCompletePartIn
	OperatorId int64
}

type MultipartAbortInput struct {
	UploadId   string
	OperatorId int64
}

type IMedia interface {
	Upload(ctx context.Context, in UploadInput) (*UploadDTO, error)
	ReadObject(ctx context.Context, rawURL, objectKey string) ([]byte, string, error)
	MultipartInit(ctx context.Context, in MultipartInitInput) (*MultipartInitDTO, error)
	MultipartPresign(ctx context.Context, in MultipartPresignInput) ([]MultipartPresignItemDTO, error)
	MultipartParts(ctx context.Context, uploadId string, operatorId int64) (*MultipartPartsDTO, error)
	MultipartComplete(ctx context.Context, in MultipartCompleteInput) (*UploadDTO, error)
	MultipartAbort(ctx context.Context, in MultipartAbortInput) error
}
