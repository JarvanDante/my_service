// Package front 前台UGC投稿控制器。
package front

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/publish/v1"
	"github.com/JarvanDante/my_service/internal/modules/publish/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IPublish }

func New(svc service.IPublish) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func toItem(d *service.PublishDTO) v1.Item {
	return v1.Item{
		Id: d.Id, UserId: d.UserId, Type: d.Type, Title: d.Title, Intro: d.Intro,
		Cover: d.Cover, Resource: d.Resource, Tags: d.Tags, Status: d.Status,
		RejectReason: d.RejectReason, AuditAt: d.AuditAt, CreatedAt: d.CreatedAt,
	}
}

func (c *Controller) Submit(ctx context.Context, req *v1.SubmitReq) (res *v1.SubmitRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.svc.Submit(ctx, service.SubmitInput{
		UserId: userId, Type: req.Type, Title: req.Title, Intro: req.Intro,
		Cover: req.Cover, Resource: req.Resource, Tags: req.Tags,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SubmitRes{Id: id}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	status := -1 // 空字符串=全部
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.My(ctx, userId, service.ListFilter{
		Status: status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toItem(d))
	}
	return res, nil
}

func (c *Controller) Cancel(ctx context.Context, req *v1.CancelReq) (res *v1.CancelRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.Cancel(ctx, userId, req.Id); err != nil {
		return nil, err
	}
	return &v1.CancelRes{}, nil
}
