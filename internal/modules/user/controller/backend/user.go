// Package backend 后台用户控制器。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

type Controller struct{ user service.IUser }

func New(svc service.IUser) *Controller { return &Controller{user: svc} }

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	dto, err := c.user.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{Id: dto.ID, Name: dto.Name, Active: dto.Active}, nil
}

func (c *Controller) Disable(ctx context.Context, req *v1.DisableReq) (res *v1.DisableRes, err error) {
	res = &v1.DisableRes{}
	err = c.user.DisableUser(ctx, req.Id)
	return
}
