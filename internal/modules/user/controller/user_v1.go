// Package controller 用户接入适配层: 把 HTTP 请求转成 service 调用。保持极薄。
package controller

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

type Controller struct {
	user service.IUser
}

func New(svc service.IUser) *Controller {
	return &Controller{user: svc}
}

func (c *Controller) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	dto, err := c.user.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{Id: dto.ID, Name: dto.Name, Active: dto.Active}, nil
}

func (c *Controller) Disable(ctx context.Context, req *v1.DisableReq) (res *v1.DisableRes, err error) {
	res = &v1.DisableRes{}
	err = c.user.DisableUser(ctx, req.Id)
	return
}
