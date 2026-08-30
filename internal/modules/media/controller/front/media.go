package front

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/media/v1"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/aesbnc"
	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/storage"
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
	name := req.Url
	if name == "" {
		name = req.ObjectKey
	}
	if !aesbnc.IsEncryptedName(name) {
		if err = c.streamPlainObject(ctx, req); err != nil {
			return nil, err
		}
		return nil, nil
	}
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

func (c *Controller) streamPlainObject(ctx context.Context, req *v1.ObjectReq) error {
	bucket, key, name, err := c.media.ResolveObjectRef(ctx, req.Url, req.ObjectKey)
	if err != nil {
		return err
	}
	client, err := storage.Get(ctx)
	if err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "对象存储不可用")
	}
	info, err := client.StatIn(ctx, bucket, key)
	if err != nil {
		return gerror.WrapCode(gcode.CodeNotFound, err, "对象不存在")
	}
	size := info.Size
	start, end, partial := parseBytesRange(ghttp.RequestFromCtx(ctx).Header.Get("Range"), size)
	length := int64(0)
	if partial {
		length = end - start + 1
	}
	obj, _, err := client.OpenIn(ctx, bucket, key, start, length)
	if err != nil {
		return gerror.WrapCode(gcode.CodeNotFound, err, "对象不存在")
	}
	defer obj.Close()

	r := ghttp.RequestFromCtx(ctx)
	ct := info.ContentType
	if ct == "" || ct == "application/octet-stream" {
		ct = sniffVideoType(name)
	}
	r.Response.Header().Set("Content-Type", ct)
	r.Response.Header().Set("Accept-Ranges", "bytes")
	r.Response.Header().Set("Cache-Control", "private, max-age=120")
	if partial {
		r.Response.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		r.Response.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
		r.Response.WriteHeader(http.StatusPartialContent)
	} else {
		r.Response.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	_, _ = io.Copy(r.Response.Writer, obj)
	r.ExitAll()
	return nil
}

func parseBytesRange(header string, size int64) (start, end int64, partial bool) {
	end = size - 1
	if size <= 0 || !strings.HasPrefix(header, "bytes=") {
		return 0, end, false
	}
	spec := strings.TrimPrefix(header, "bytes=")
	if i := strings.IndexByte(spec, ','); i >= 0 {
		spec = spec[:i]
	}
	parts := strings.SplitN(spec, "-", 2)
	if len(parts) != 2 {
		return 0, end, false
	}
	if parts[0] == "" {
		n, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || n <= 0 {
			return 0, end, false
		}
		if n > size {
			n = size
		}
		return size - n, end, true
	}
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || n < 0 || n >= size {
		return 0, end, false
	}
	start = n
	if parts[1] != "" {
		n, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || n < start {
			return 0, end, false
		}
		if n < size {
			end = n
		}
	}
	return start, end, true
}

func sniffVideoType(name string) string {
	low := strings.ToLower(name)
	switch {
	case strings.Contains(low, ".webm"):
		return "video/webm"
	case strings.Contains(low, ".mov"):
		return "video/quicktime"
	case strings.Contains(low, ".mkv"):
		return "video/x-matroska"
	default:
		return "video/mp4"
	}
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
