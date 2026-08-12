// Package front 前台基础配置控制器(公开)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/config/v1"
	"github.com/JarvanDante/my_service/internal/modules/config/service"
)

type Controller struct{ svc service.IConfig }

func New(svc service.IConfig) *Controller { return &Controller{svc: svc} }

func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	configs, err := c.svc.Info(ctx, req.Grp)
	if err != nil {
		return nil, err
	}
	return &v1.InfoRes{Configs: configs}, nil
}

func (c *Controller) Check(ctx context.Context, req *v1.CheckReq) (res *v1.CheckRes, err error) {
	return &v1.CheckRes{}, nil
}
