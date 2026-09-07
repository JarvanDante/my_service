package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/channel/v1"
	"github.com/JarvanDante/my_service/internal/modules/channel/service"
)

type Controller struct{ svc service.IChannel }

func New(svc service.IChannel) *Controller { return &Controller{svc: svc} }

func toItem(r *service.ItemDTO) v1.Item {
	return v1.Item{
		Id:             r.Id,
		Code:           r.Code,
		Name:           r.Name,
		Remark:         r.Remark,
		Status:         r.Status,
		StatusText:     r.StatusText,
		UserCount:      r.UserCount,
		H5Link:         r.H5Link,
		ClipboardValue: r.ClipboardValue,
		CreatedAt:      r.CreatedAt,
	}
}

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
	for _, r := range list {
		res.List = append(res.List, toItem(r))
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Code: req.Code, Name: req.Name, Remark: req.Remark, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Name: req.Name, Remark: req.Remark, Status: req.Status,
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
