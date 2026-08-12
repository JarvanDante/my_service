// Package backend 后台推广应用控制器(CRUD)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/application/v1"
	"github.com/JarvanDante/my_service/internal/modules/application/service"
)

type Controller struct{ svc service.IApplication }

func New(svc service.IApplication) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	// status 空字符串 = 不过滤(全部); "0"=只看下架; "1"=只看上架。
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
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Name: r.Name, Tag: r.Tag, Intro: r.Intro, Avatar: r.Avatar,
			DownloadUrl: r.DownloadUrl, IosUrl: r.IosUrl, AndroidUrl: r.AndroidUrl,
			LocIds: r.LocIds, Rank: r.Rank, DownTotal: r.DownTotal, Status: r.Status,
			CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Name: req.Name, Tag: req.Tag, Intro: req.Intro, Avatar: req.Avatar,
		DownloadUrl: req.DownloadUrl, IosUrl: req.IosUrl, AndroidUrl: req.AndroidUrl,
		LocIds: req.LocIds, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Name: req.Name, Tag: req.Tag, Intro: req.Intro, Avatar: req.Avatar,
		DownloadUrl: req.DownloadUrl, IosUrl: req.IosUrl, AndroidUrl: req.AndroidUrl,
		LocIds: req.LocIds, Rank: req.Rank, Status: req.Status,
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
