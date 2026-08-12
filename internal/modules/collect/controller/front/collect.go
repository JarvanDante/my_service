// Package front 前台收藏/点赞控制器(需登录)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/collect/v1"
	"github.com/JarvanDante/my_service/internal/modules/collect/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.ICollect }

func New(svc service.ICollect) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Operate(ctx context.Context, req *v1.OperateReq) (res *v1.OperateRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	ids := req.Ids
	if len(ids) == 0 && req.Id > 0 {
		ids = []int64{req.Id}
	}
	if err = c.svc.Operate(ctx, service.OperateInput{
		UserId: userId, Ids: ids, MediaType: req.MediaType, Flag: req.Flag, Type: req.Type,
	}); err != nil {
		return nil, err
	}
	return &v1.OperateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.Delete(ctx, userId, req.Ids, req.MediaType, req.Type); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		UserId: userId, Type: req.Type, MediaType: req.MediaType,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.CollectItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.CollectItem{
			ContentId: r.ContentId, MediaType: r.MediaType, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}
