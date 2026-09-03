package front

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/captcha/v1"
	"github.com/JarvanDante/my_service/internal/modules/captcha/service"
)

type Controller struct{ svc service.ICaptcha }

func New(svc service.ICaptcha) *Controller { return &Controller{svc: svc} }

func (c *Controller) Issue(ctx context.Context, req *v1.IssueReq) (res *v1.IssueRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	id, image, err := c.svc.Issue(ctx, r.GetClientIp())
	if err != nil {
		return nil, err
	}
	return &v1.IssueRes{Id: id, Image: image}, nil
}

func (c *Controller) Verify(ctx context.Context, req *v1.VerifyReq) (res *v1.VerifyRes, err error) {
	if err = c.svc.Verify(ctx, req.Id, req.Code); err != nil {
		return nil, err
	}
	return &v1.VerifyRes{Ok: true}, nil
}
