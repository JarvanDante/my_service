// Package backend 后台基础配置控制器(KV 管理)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/config/v1"
	"github.com/JarvanDante/my_service/internal/modules/config/service"
)

type Controller struct{ svc service.IConfig }

func New(svc service.IConfig) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	// status 空字符串 = 不过滤(全部); "0"=只看禁用; "1"=只看启用。
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Grp: req.Grp, Status: statusFilter, Keyword: req.Keyword,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Grp: r.Grp, Key: r.Key, Value: r.Value,
			Remark: r.Remark, Status: r.Status, UpdatedAt: r.UpdatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.CreateInput{
		Grp: req.Grp, Key: req.Key, Value: req.Value,
		Remark: req.Remark, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.UpdateInput{
		Id: req.Id, Grp: req.Grp, Value: req.Value,
		Remark: req.Remark, Status: req.Status,
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
