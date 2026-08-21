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
	case "image", "avatar", "video":
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "前台仅支持上传 image/avatar/video")
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
