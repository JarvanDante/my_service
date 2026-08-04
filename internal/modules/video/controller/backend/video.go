package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/video/v1"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ video service.IVideo }

func New(svc service.IVideo) *Controller { return &Controller{video: svc} }

func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.video.List(ctx, service.ListInput{
		Keyword: req.Keyword, Status: req.Status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.VideoItem, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.ListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.video.Create(ctx, service.SaveInput{
		Title: req.Title, Description: req.Description,
		CoverUrl: req.CoverUrl, CoverKey: req.CoverKey, CoverMediaId: req.CoverMediaId,
		SourceUrl: req.SourceUrl, SourceKey: req.SourceKey, SourceMediaId: req.SourceMediaId,
		Duration: req.Duration, Sort: req.Sort, Status: req.Status, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	err = c.video.Update(ctx, service.SaveInput{
		Id: req.Id, Title: req.Title, Description: req.Description,
		CoverUrl: req.CoverUrl, CoverKey: req.CoverKey, CoverMediaId: req.CoverMediaId,
		SourceUrl: req.SourceUrl, SourceKey: req.SourceKey, SourceMediaId: req.SourceMediaId,
		Duration: req.Duration, Sort: req.Sort, Status: req.Status, OperatorId: op,
	})
	return &v1.UpdateRes{}, err
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	return &v1.DeleteRes{}, c.video.Delete(ctx, req.Id)
}

func (c *Controller) Status(ctx context.Context, req *v1.StatusReq) (res *v1.StatusRes, err error) {
	return &v1.StatusRes{}, c.video.SetStatus(ctx, req.Id, req.Status)
}

func toItem(v *service.VideoDTO) v1.VideoItem {
	if v == nil {
		return v1.VideoItem{}
	}
	return v1.VideoItem{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, CoverKey: v.CoverKey, CoverMediaId: v.CoverMediaId,
		SourceUrl: v.SourceUrl, SourceKey: v.SourceKey, SourceMediaId: v.SourceMediaId,
		Duration: v.Duration, Sort: v.Sort, Status: v.Status, CreatedBy: v.CreatedBy,
		CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}
