// Package backend 后台图集控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/photo/v1"
	"github.com/JarvanDante/my_service/internal/modules/photo/service"
)

type Controller struct{ svc service.IPhoto }

func New(svc service.IPhoto) *Controller { return &Controller{svc: svc} }

// toPicDTO 保留 nil 语义: 编辑时"没传 pics"与"传了空数组(清空)"是两回事。
func toPicDTO(in []v1.Pic) []service.PicDTO {
	if in == nil {
		return nil
	}
	out := make([]service.PicDTO, 0, len(in))
	for _, p := range in {
		out = append(out, service.PicDTO{Url: p.Url, Width: p.Width, Height: p.Height})
	}
	return out
}

func toPicV1(in []service.PicDTO) []v1.Pic {
	out := make([]v1.Pic, 0, len(in))
	for _, p := range in {
		out = append(out, v1.Pic{Url: p.Url, Width: p.Width, Height: p.Height})
	}
	return out
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	status := -1 // 空字符串=不筛选, 与其它后台列表口径一致
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: status, Category: req.Category, Keyword: req.Keyword,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, Title: d.Title, Cover: d.Cover, Intro: d.Intro, Category: d.Category,
			Tags: d.Tags, IsVip: d.IsVip, Price: d.Price, FreeCount: d.FreeCount,
			Pics: toPicV1(d.Pics), PicCount: d.PicCount, ViewCount: d.ViewCount,
			BuyCount: d.BuyCount, LikeCount: d.LikeCount, Rank: d.Rank, Status: d.Status,
			PublishId: d.PublishId, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Title: req.Title, Cover: req.Cover, Intro: req.Intro, Category: req.Category,
		Tags: req.Tags, IsVip: req.IsVip, Price: req.Price, FreeCount: req.FreeCount,
		Pics: toPicDTO(req.Pics), Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Title: req.Title, Cover: req.Cover, Intro: req.Intro,
		Category: req.Category, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price,
		FreeCount: req.FreeCount, Pics: toPicDTO(req.Pics), Rank: req.Rank, Status: req.Status,
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

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}
