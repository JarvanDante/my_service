// Package backend 后台系统消息控制器(发布/管理)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/message/v1"
	"github.com/JarvanDante/my_service/internal/modules/message/service"
)

type Controller struct{ svc service.IMessage }

func New(svc service.IMessage) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	// status 空字符串 = 不过滤(全部); "0"=撤回; "1"=发布。
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	// user_id 空字符串 = 不过滤(全部); "0"=只看全员; ">0"=指定用户。
	var userFilter int64 = -1
	if req.UserId != "" {
		if v, e := strconv.ParseInt(req.UserId, 10, 64); e == nil {
			userFilter = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		UserId: userFilter, Status: statusFilter, Keyword: req.Keyword,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, UserId: r.UserId, Type: r.Type, Title: r.Title,
			Content: r.Content, Status: r.Status, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.CreateInput{
		UserId: req.UserId, Type: req.Type, Title: req.Title, Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.UpdateInput{
		Id: req.Id, Type: req.Type, Title: req.Title, Content: req.Content, Status: req.Status,
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
