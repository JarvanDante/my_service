// Package backend 后台商品兑换控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/redeemgoods/v1"
	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/service"
)

type Controller struct{ svc service.IRedeemGoods }

func New(svc service.IRedeemGoods) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: statusFilter, Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, Name: d.Name, Cover: d.Cover, Intro: d.Intro, CostGold: d.CostGold,
			Stock: d.Stock, Exchanged: d.Exchanged, Rank: d.Rank, Status: d.Status,
			CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Name: req.Name, Cover: req.Cover, Intro: req.Intro, CostGold: req.CostGold,
		Stock: req.Stock, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Name: req.Name, Cover: req.Cover, Intro: req.Intro,
		CostGold: req.CostGold, Stock: req.Stock, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err = c.svc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) Orders(ctx context.Context, req *v1.OrdersReq) (res *v1.OrdersRes, err error) {
	list, total, err := c.svc.Orders(ctx, service.ListFilter{
		UserId: req.UserId, GoodsId: req.GoodsId, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.OrdersRes{Total: total, List: make([]v1.OrderItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.OrderItem{
			Id: d.Id, UserId: d.UserId, GoodsId: d.GoodsId, GoodsName: d.GoodsName,
			CostGold: d.CostGold, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}
