// Package backend 后台优惠券控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/coupon/v1"
	"github.com/JarvanDante/my_service/internal/modules/coupon/service"
)

type Controller struct{ svc service.ICoupon }

func New(svc service.ICoupon) *Controller { return &Controller{svc: svc} }

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, e := strconv.Atoi(s); e == nil {
		return v
	}
	return def
}

func atoi64(s string) int64 {
	if v, e := strconv.ParseInt(s, 10, 64); e == nil {
		return v
	}
	return 0
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: atoiDefault(req.Status, -1), Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, Name: d.Name, Type: d.Type, Scene: d.Scene, FaceValue: d.FaceValue,
			Discount: d.Discount, Threshold: d.Threshold, MaxDeduct: d.MaxDeduct,
			Total: d.Total, Issued: d.Issued, PerLimit: d.PerLimit, ExpireDay: d.ExpireDay,
			Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.TplInput{
		Name: req.Name, Type: req.Type, Scene: req.Scene, FaceValue: req.FaceValue,
		Discount: req.Discount, Threshold: req.Threshold, MaxDeduct: req.MaxDeduct,
		Total: req.Total, PerLimit: req.PerLimit, ExpireDay: req.ExpireDay, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.TplInput{
		Id: req.Id, Name: req.Name, Type: req.Type, Scene: req.Scene, FaceValue: req.FaceValue,
		Discount: req.Discount, Threshold: req.Threshold, MaxDeduct: req.MaxDeduct,
		Total: req.Total, PerLimit: req.PerLimit, ExpireDay: req.ExpireDay, Status: req.Status,
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

func (c *Controller) Grant(ctx context.Context, req *v1.GrantReq) (res *v1.GrantRes, err error) {
	ok, fail, errs := c.svc.Grant(ctx, req.Id, req.UserIds)
	return &v1.GrantRes{Success: ok, Failed: fail, Errors: errs}, nil
}

func (c *Controller) Users(ctx context.Context, req *v1.UsersReq) (res *v1.UsersRes, err error) {
	list, total, err := c.svc.Users(ctx, service.ListFilter{
		TplId: atoi64(req.TplId), UserId: atoi64(req.UserId),
		Status: atoiDefault(req.Status, 0), Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.UsersRes{Total: total, List: make([]v1.UserItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.UserItem{
			Id: d.Id, UserId: d.UserId, TplId: d.TplId, Name: d.Name, Type: d.Type,
			FaceValue: d.FaceValue, Discount: d.Discount, Status: d.Status,
			StatusText: service.StatusText(d.Status), RefId: d.RefId,
			ExpireAt: d.ExpireAt, UsedAt: d.UsedAt, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}
