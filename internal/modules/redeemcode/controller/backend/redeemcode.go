// Package backend 后台兑换码控制器(建码/管理/使用记录)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/redeemcode/v1"
	"github.com/JarvanDante/my_service/internal/modules/redeemcode/service"
)

type Controller struct{ svc service.IRedeemCode }

func New(svc service.IRedeemCode) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	// status 空字符串 = 不过滤(全部); "0"=只看禁用; "1"=只看启用。
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
			Id: r.Id, Name: r.Name, Code: r.Code, CardType: r.CardType, Value: r.Value,
			TotalTimes: r.TotalTimes, UsedTimes: r.UsedTimes, Status: r.Status,
			ExpiredAt: r.ExpiredAt, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, code, err := c.svc.Create(ctx, service.CreateInput{
		Name: req.Name, Code: req.Code, Value: req.Value,
		TotalTimes: req.TotalTimes, ExpiredAt: req.ExpiredAt, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id, Code: code}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.UpdateInput{
		Id: req.Id, Name: req.Name, Value: req.Value,
		TotalTimes: req.TotalTimes, ExpiredAt: req.ExpiredAt, Status: req.Status,
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

func (c *Controller) Records(ctx context.Context, req *v1.RecordListReq) (res *v1.RecordListRes, err error) {
	list, total, err := c.svc.Records(ctx, service.RecordFilter{
		UserId: req.UserId, Code: req.Code, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.RecordListRes{Total: total, List: make([]v1.RecordItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.RecordItem{
			Id: r.Id, UserId: r.UserId, CodeId: r.CodeId, Code: r.Code, Name: r.Name,
			CardType: r.CardType, Value: r.Value, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}
