package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/media/v1"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/aesbnc"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ media service.IMedia }

func New(svc service.IMedia) *Controller { return &Controller{media: svc} }

func (c *Controller) Preview(ctx context.Context, req *v1.PreviewReq) (res *v1.PreviewRes, err error) {
	data, name, err := c.media.ReadObject(ctx, req.Url, req.ObjectKey)
	if err != nil {
		return nil, err
	}
	out := data
	if aesbnc.IsEncryptedName(name) || !aesbnc.LooksLikeImage(data) {
		dec, derr := aesbnc.Decrypt(data)
		if derr == nil {
			out = dec
		}
	}
	ct := aesbnc.SniffContentType(out)
	if ct == "application/octet-stream" {
		ct = "image/jpeg"
	}
	r := ghttp.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", ct)
	r.Response.Header().Set("Cache-Control", "private, max-age=120")
	r.Response.Write(out)
	r.ExitAll()
	return nil, nil
}

func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

// Upload 整文件上传(小文件/封面)。
func (c *Controller) Upload(ctx context.Context, req *v1.UploadReq) (res *v1.UploadRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.media.Upload(ctx, service.UploadInput{
		File: req.File, Purpose: req.Purpose, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	return toUploadRes(dto), nil
}

// MultipartInit 初始化分片上传。
func (c *Controller) MultipartInit(ctx context.Context, req *v1.MultipartInitReq) (res *v1.MultipartInitRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.media.MultipartInit(ctx, service.MultipartInitInput{
		Filename: req.Filename, Purpose: req.Purpose, ContentType: req.ContentType,
		Size: req.Size, PartSize: req.PartSize, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MultipartInitRes{
		UploadId: dto.UploadId, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
		PartSize: dto.PartSize, PartCount: dto.PartCount,
	}, nil
}

// MultipartPresign 批量预签名分片 PUT。
func (c *Controller) MultipartPresign(ctx context.Context, req *v1.MultipartPresignReq) (res *v1.MultipartPresignRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	list, err := c.media.MultipartPresign(ctx, service.MultipartPresignInput{
		UploadId: req.UploadId, PartNumbers: req.PartNumbers, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.MultipartPresignItem, 0, len(list))
	for _, it := range list {
		items = append(items, v1.MultipartPresignItem{
			PartNumber: it.PartNumber, Url: it.Url, Method: it.Method, ExpiresIn: it.ExpiresIn,
		})
	}
	return &v1.MultipartPresignRes{List: items}, nil
}

// MultipartParts 已上传分片列表(断点续传)。
func (c *Controller) MultipartParts(ctx context.Context, req *v1.MultipartPartsReq) (res *v1.MultipartPartsRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.media.MultipartParts(ctx, req.UploadId, op)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MultipartPartItem, 0, len(dto.List))
	for _, p := range dto.List {
		items = append(items, v1.MultipartPartItem{
			PartNumber: p.PartNumber, Etag: p.Etag, Size: p.Size,
		})
	}
	return &v1.MultipartPartsRes{
		UploadId: dto.UploadId, Status: dto.Status, PartCount: dto.PartCount, List: items,
	}, nil
}

// MultipartComplete 合并分片。
func (c *Controller) MultipartComplete(ctx context.Context, req *v1.MultipartCompleteReq) (res *v1.MultipartCompleteRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	parts := make([]service.MultipartCompletePartIn, 0, len(req.Parts))
	for _, p := range req.Parts {
		parts = append(parts, service.MultipartCompletePartIn{PartNumber: p.PartNumber, Etag: p.Etag})
	}
	dto, err := c.media.MultipartComplete(ctx, service.MultipartCompleteInput{
		UploadId: req.UploadId, Parts: parts, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MultipartCompleteRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
	}, nil
}

// MultipartAbort 取消分片上传。
func (c *Controller) MultipartAbort(ctx context.Context, req *v1.MultipartAbortReq) (res *v1.MultipartAbortRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.media.MultipartAbort(ctx, service.MultipartAbortInput{
		UploadId: req.UploadId, OperatorId: op,
	}); err != nil {
		return nil, err
	}
	return &v1.MultipartAbortRes{}, nil
}

func toUploadRes(dto *service.UploadDTO) *v1.UploadRes {
	return &v1.UploadRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
	}
}
