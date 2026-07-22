// Package front 前台用户控制器。ghttp 适配层, 极薄。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

type Controller struct{ user service.IUser }

func New(svc service.IUser) *Controller { return &Controller{user: svc} }

func (c *Controller) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	dto, err := c.user.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{Id: dto.ID, Name: dto.Name, Active: dto.Active}, nil
}
