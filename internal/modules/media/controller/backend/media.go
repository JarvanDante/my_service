package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/media/v1"
	"github.com/JarvanDante/my_service/internal/modules/media/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ media service.IMedia }

func New(svc service.IMedia) *Controller { return &Controller{media: svc} }

func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

// Upload 上传媒体到 MinIO。
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
	return &v1.UploadRes{
		Id: dto.Id, Url: dto.Url, ObjectKey: dto.ObjectKey, Bucket: dto.Bucket,
		Purpose: dto.Purpose, ContentType: dto.ContentType, Size: dto.Size,
	}, nil
}
