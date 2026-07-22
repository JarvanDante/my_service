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

// Login 公开接口。
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

// Info 需登录: userId 由 Auth 中间件写入 ctx。
func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	uid := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if uid <= 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	dto, err := c.user.Info(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &v1.InfoRes{UserInfo: toApiUser(dto)}, nil
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
