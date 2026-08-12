// Package front 前台推广应用控制器(公开)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/application/v1"
	"github.com/JarvanDante/my_service/internal/modules/application/service"
)

type Controller struct{ svc service.IApplication }

func New(svc service.IApplication) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := c.svc.FrontList(ctx, req.Loc)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.AppItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.AppItem{
			Id: r.Id, Name: r.Name, Tag: r.Tag, Desc: r.Intro, Avatar: r.Avatar,
			DownloadUrl: r.DownloadUrl, IosUrl: r.IosUrl, AndroidUrl: r.AndroidUrl,
			LocIds: r.LocIds, Rank: r.Rank, DownTotal: r.DownTotal,
		})
	}
	return res, nil
}

func (c *Controller) Click(ctx context.Context, req *v1.ClickReq) (res *v1.ClickRes, err error) {
	if err = c.svc.Click(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ClickRes{}, nil
}
