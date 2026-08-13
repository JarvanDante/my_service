// Package front 前台钱包控制器。
package front

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/wallet/v1"
	"github.com/JarvanDante/my_service/internal/modules/wallet/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IWallet }

func New(svc service.IWallet) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Balance(ctx context.Context, req *v1.BalanceReq) (res *v1.BalanceRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Summary(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.BalanceRes{
		Balance: d.Balance, Frozen: d.Frozen, TotalIn: d.TotalIn,
		TotalOut: d.TotalOut, Withdrawn: d.Withdrawn,
	}, nil
}

func (c *Controller) Waters(ctx context.Context, req *v1.WatersReq) (res *v1.WatersRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	direction := -1
	if req.Direction != "" {
		if v, e := strconv.Atoi(req.Direction); e == nil {
			direction = v
		}
	}
	list, total, err := c.svc.Waters(ctx, service.WaterFilter{
		UserId: userId, Direction: direction, Scene: req.Scene, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.WatersRes{Total: total, List: make([]v1.WaterItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.WaterItem{
			Id: d.Id, Direction: d.Direction, Scene: d.Scene, Amount: d.Amount,
			BalanceBefore: d.BalanceBefore, BalanceAfter: d.BalanceAfter,
			Remark: d.Remark, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}
