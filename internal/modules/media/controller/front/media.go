package front

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/media/v1"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/aesbnc"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ media service.IMedia }

func New(svc service.IMedia) *Controller { return &Controller{media: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Object(ctx context.Context, req *v1.ObjectReq) (res *v1.ObjectRes, err error) {
	data, name, err := c.media.ReadObject(ctx, req.Url, req.ObjectKey)
	if err != nil {
		return nil, err
	}
	r := ghttp.RequestFromCtx(ctx)
	ct := "application/octet-stream"
	if !aesbnc.IsEncryptedName(name) && aesbnc.LooksLikeImage(data) {
		ct = aesbnc.SniffContentType(data)
	}
	r.Response.Header().Set("Content-Type", ct)
	r.Response.Header().Set("Cache-Control", "private, max-age=120")
	r.Response.Write(data)
	r.ExitAll()
	return nil, nil
}

func (c *Controller) Upload(ctx context.Context, req *v1.UploadReq) (res *v1.UploadRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	purpose := strings.ToLower(strings.TrimSpace(req.Purpose))
	if purpose == "" {
		purpose = "image"
	}
	switch purpose {
	case "image", "avatar", "video", "ad", "post", "post_video":
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "前台仅支持上传 image/avatar/video/ad/post")
	}
	dto, err := c.media.Upload(ctx, service.UploadInput{
		File: req.File, Purpose: purpose, OperatorId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UploadRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
	}, nil
}

func (c *Controller) StorageInit(ctx context.Context, req *v1.StorageInitReq) (res *v1.StorageInitRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	dto, err := c.media.InitStorageUpload(ctx, service.StorageInitInput{
		Filename: req.Filename, Purpose: req.Purpose,
		ContentType: req.ContentType, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &v1.StorageInitRes{
		Id: dto.Id, UploadUrl: dto.UploadUrl, Method: dto.Method, Bucket: dto.Bucket,
		ObjectKey: dto.ObjectKey, ExpireSec: dto.ExpireSec, PublicUrl: dto.PublicUrl,
		ContentType: dto.ContentType,
	}, nil
}

func (c *Controller) StorageConfirm(ctx context.Context, req *v1.StorageConfirmReq) (res *v1.StorageConfirmRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	dto, err := c.media.ConfirmStorageUpload(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.StorageConfirmRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket, Size: dto.Size,
	}, nil
}

func (c *Controller) MultipartInit(ctx context.Context, req *v1.MultipartInitReq) (res *v1.MultipartInitRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	return nil, gerror.NewCode(gcode.CodeInvalidParameter, "前台视频请走统一存储，不再分片上传到 my-media")
}

func (c *Controller) MultipartPart(ctx context.Context, req *v1.MultipartPartUploadReq) (res *v1.MultipartPartUploadRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.media.MultipartUploadPart(ctx, service.MultipartUploadPartInput{
		UploadId: req.UploadId, PartNumber: req.PartNumber, File: req.File, OperatorId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MultipartPartUploadRes{PartNumber: dto.PartNumber, Etag: dto.Etag, Size: dto.Size}, nil
}

func (c *Controller) MultipartParts(ctx context.Context, req *v1.MultipartPartsReq) (res *v1.MultipartPartsRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.media.MultipartParts(ctx, req.UploadId, userId)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MultipartPartItem, 0, len(dto.List))
	for _, p := range dto.List {
		items = append(items, v1.MultipartPartItem{PartNumber: p.PartNumber, Etag: p.Etag, Size: p.Size})
	}
	return &v1.MultipartPartsRes{
		UploadId: dto.UploadId, Status: dto.Status, PartCount: dto.PartCount, List: items,
	}, nil
}

func (c *Controller) MultipartComplete(ctx context.Context, req *v1.MultipartCompleteReq) (res *v1.MultipartCompleteRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	parts := make([]service.MultipartCompletePartIn, 0, len(req.Parts))
	for _, p := range req.Parts {
		parts = append(parts, service.MultipartCompletePartIn{PartNumber: p.PartNumber, Etag: p.Etag})
	}
	dto, err := c.media.MultipartComplete(ctx, service.MultipartCompleteInput{
		UploadId: req.UploadId, Parts: parts, OperatorId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MultipartCompleteRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
	}, nil
}

func (c *Controller) MultipartAbort(ctx context.Context, req *v1.MultipartAbortReq) (res *v1.MultipartAbortRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.media.MultipartAbort(ctx, service.MultipartAbortInput{
		UploadId: req.UploadId, OperatorId: userId,
	}); err != nil {
		return nil, err
	}
	return &v1.MultipartAbortRes{}, nil
}
