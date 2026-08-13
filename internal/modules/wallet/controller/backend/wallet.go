// Package backend 后台钱包控制器。
package backend

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/wallet/v1"
	"github.com/JarvanDante/my_service/internal/modules/wallet/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IWallet }

func New(svc service.IWallet) *Controller { return &Controller{svc: svc} }

func adminId(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
}

func (c *Controller) Logs(ctx context.Context, req *v1.LogsReq) (res *v1.LogsRes, err error) {
	var userId int64
	if req.UserId != "" {
		if v, e := strconv.ParseInt(req.UserId, 10, 64); e == nil {
			userId = v
		}
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
	res = &v1.LogsRes{Total: total, List: make([]v1.LogItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.LogItem{
			Id: d.Id, UserId: d.UserId, Direction: d.Direction, Scene: d.Scene,
			Amount: d.Amount, BalanceBefore: d.BalanceBefore, BalanceAfter: d.BalanceAfter,
			RefId: d.RefId, Remark: d.Remark, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Adjust(ctx context.Context, req *v1.AdjustReq) (res *v1.AdjustRes, err error) {
	bal, err := c.svc.Adjust(ctx, adminId(ctx), req.UserId, req.Amount, req.Remark)
	if err != nil {
		return nil, err
	}
	return &v1.AdjustRes{Balance: bal}, nil
}
