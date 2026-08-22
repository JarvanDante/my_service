// Package front 前台意见反馈控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/feedback/v1"
	"github.com/JarvanDante/my_service/internal/modules/feedback/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IFeedback }

func New(svc service.IFeedback) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Add(ctx context.Context, req *v1.AddReq) (res *v1.AddRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	newId, err := c.svc.Add(ctx, service.AddInput{
		UserId: id, Type: req.Type, ProblemType: req.ProblemType, Content: req.Content,
		Pics: req.Pics, SysInfo: req.SysInfo, MediaId: req.MediaId, MediaTitle: req.MediaTitle,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AddRes{Id: newId}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		UserId: userId, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.MyItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.MyItem{
			Id: r.Id, ProblemType: r.ProblemType, Content: r.Content,
			Status: r.Status, Reply: r.Reply, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}
