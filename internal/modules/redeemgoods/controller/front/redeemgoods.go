// Package front 前台商品兑换控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/redeemgoods/v1"
	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IRedeemGoods }

func New(svc service.IRedeemGoods) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.FrontList(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.GoodsItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.GoodsItem{
			Id: d.Id, Name: d.Name, Cover: d.Cover, Intro: d.Intro,
			CostGold: d.CostGold, Stock: d.Stock,
		})
	}
	return res, nil
}

func (c *Controller) Exchange(ctx context.Context, req *v1.ExchangeReq) (res *v1.ExchangeRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	orderId, err := c.svc.Exchange(ctx, userId, req.GoodsId)
	if err != nil {
		return nil, err
	}
	return &v1.ExchangeRes{OrderId: orderId}, nil
}

func (c *Controller) History(ctx context.Context, req *v1.HistoryReq) (res *v1.HistoryRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, err := c.svc.History(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.HistoryRes{List: make([]v1.HistoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.HistoryItem{
			Id: d.Id, GoodsName: d.GoodsName, CostGold: d.CostGold, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}
