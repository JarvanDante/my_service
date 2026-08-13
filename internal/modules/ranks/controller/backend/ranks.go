// Package backend 后台热搜/排行控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/ranks/v1"
	"github.com/JarvanDante/my_service/internal/modules/ranks/service"
)

type Controller struct{ svc service.IRank }

func New(svc service.IRank) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.HotList(ctx, service.HotFilter{
		Status: statusFilter, Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Keyword: r.Keyword, Heat: r.Heat,
			SearchCount: r.SearchCount, Status: r.Status, UpdatedAt: r.UpdatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.HotCreate(ctx, req.Keyword, req.Heat, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.HotUpdate(ctx, req.Id, req.Keyword, req.Heat, req.Status); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err = c.svc.HotDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) RefreshRank(ctx context.Context, req *v1.RefreshRankReq) (res *v1.RefreshRankRes, err error) {
	if err = c.svc.RefreshRank(ctx); err != nil {
		return nil, err
	}
	return &v1.RefreshRankRes{}, nil
}
