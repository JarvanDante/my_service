// Package backend 后台帖子控制器(审核/管理)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/post/v1"
	"github.com/JarvanDante/my_service/internal/modules/post/service"
)

type Controller struct {
	svc service.IPost
	cat service.ICategory
}

func New(svc service.IPost, cat service.ICategory) *Controller {
	return &Controller{svc: svc, cat: cat}
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: statusFilter, Keyword: req.Keyword, UserId: req.UserId,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, UserId: d.UserId, Title: d.Title, Content: d.Content, Pics: d.Pics,
			Topics: d.Topics, Category: d.Category, VideoUrl: d.VideoUrl,
			MediaId: d.MediaId, ViewCount: d.ViewCount, Rank: d.Rank, LikeCount: d.LikeCount,
			CommentCount: d.CommentCount, Status: d.Status, RejectReason: d.RejectReason,
			CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.UpdateInput{
		Id: req.Id, Category: req.Category, ViewCount: req.ViewCount, Rank: req.Rank,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, req.Id, req.Pass, req.Reason); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err = c.svc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
