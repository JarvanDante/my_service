// Package front 前台用户控制器(适配层, 薄)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ user service.IUser }

func New(svc service.IUser) *Controller { return &Controller{user: svc} }

// uid 从 ctx 取当前用户(Auth 中间件写入)。
func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

// Login 公开。
func (c *Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	dto, err := c.user.Login(ctx, service.LoginInput{
		DeviceId:      req.DeviceId,
		DeviceType:    req.DeviceType,
		DeviceVersion: req.DeviceVersion,
		Ip:            r.GetClientIp(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{Token: dto.Token, User: toApiUser(dto.User)}, nil
}

// Info 需登录。
func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.Info(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.InfoRes{UserInfo: toApiUser(dto)}, nil
}

// Logout 退出登录。
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.Logout(ctx, id); err != nil {
		return nil, err
	}
	return &v1.LogoutRes{}, nil
}

// Refresh 刷新 token。
func (c *Controller) Refresh(ctx context.Context, req *v1.RefreshReq) (res *v1.RefreshRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	token, err := c.user.Refresh(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.RefreshRes{Token: token}, nil
}

// BindPhone 绑定手机。
func (c *Controller) BindPhone(ctx context.Context, req *v1.BindPhoneReq) (res *v1.BindPhoneRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.BindPhone(ctx, id, req.Phone, req.Code); err != nil {
		return nil, err
	}
	return &v1.BindPhoneRes{}, nil
}

func toApiUser(d *service.UserInfoDTO) v1.UserInfo {
	if d == nil {
		return v1.UserInfo{}
	}
	return v1.UserInfo{
		Id: d.Id, Username: d.Username, Nickname: d.Nickname, Phone: d.Phone,
		Img: d.Img, Signature: d.Signature, Sex: d.Sex, Level: d.Level,
		Balance: d.Balance, Credit: d.Credit, GroupName: d.GroupName,
		Fans: d.Fans, Follow: d.Follow,
	}
}
