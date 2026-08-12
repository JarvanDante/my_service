// Package front 前台兑换码控制器(需登录)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/redeemcode/v1"
	"github.com/JarvanDante/my_service/internal/modules/redeemcode/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IRedeemCode }

func New(svc service.IRedeemCode) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Use(ctx context.Context, req *v1.UseReq) (res *v1.UseRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.svc.Use(ctx, userId, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.UseRes{Desc: r.Desc}, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, err := c.svc.MyRecords(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.RecordItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.RecordItem{
			Code: r.Code, Desc: r.Desc, ActivedAt: r.ActivedAt,
		})
	}
	return res, nil
}
