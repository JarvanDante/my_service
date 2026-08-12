// Package backend 后台意见反馈控制器。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/feedback/v1"
	"github.com/JarvanDante/my_service/internal/modules/feedback/service"
)

type Controller struct{ svc service.IFeedback }

func New(svc service.IFeedback) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Page: req.Page, Size: req.Size, Status: req.Status, Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, UserId: r.UserId, Type: r.Type, ProblemType: r.ProblemType,
			Content: r.Content, Pics: r.Pics, SysInfo: r.SysInfo, MediaId: r.MediaId,
			MediaTitle: r.MediaTitle, Status: r.Status, Reply: r.Reply, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Handle(ctx context.Context, req *v1.HandleReq) (res *v1.HandleRes, err error) {
	if err = c.svc.Handle(ctx, req.Id, req.Reply, req.Status); err != nil {
		return nil, err
	}
	return &v1.HandleRes{}, nil
}
