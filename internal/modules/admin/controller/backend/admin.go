// Package backend 后台管理员控制器。
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/genv"

	v1 "github.com/JarvanDante/my_service/api/backend/admin/v1"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/siteconf"
)

type Controller struct{ admin service.IAdmin }

func New(svc service.IAdmin) *Controller { return &Controller{admin: svc} }

func adminId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

// Login 公开。
func (c *Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	dto, err := c.admin.Login(ctx, service.LoginInput{Username: req.Username, Password: req.Password, Ip: r.GetClientIp()})
	if err != nil {
		return nil, err
	}
	code, name := siteMeta(ctx)
	return &v1.LoginRes{Token: dto.Token, Admin: toApi(dto.Admin), SiteCode: code, SiteName: name}, nil
}

// Logout 需登录。
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	id, err := adminId(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.admin.Logout(ctx, id); err != nil {
		return nil, err
	}
	return &v1.LogoutRes{}, nil
}

// Info 需登录。
func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	id, err := adminId(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.admin.Info(ctx, id)
	if err != nil {
		return nil, err
	}
	code, name := siteMeta(ctx)
	return &v1.InfoRes{AdminInfo: toApi(dto), SiteCode: code, SiteName: name}, nil
}

func toApi(d *service.AdminInfoDTO) v1.AdminInfo {
	if d == nil {
		return v1.AdminInfo{}
	}
	return v1.AdminInfo{Id: d.Id, Username: d.Username, Nickname: d.Nickname, RoleId: d.RoleId}
}

// siteMeta 站点品牌: SITE_CODE 来自部署 env; 站点名优先 site_config 表(后台可改), 兜底配置文件。
func siteMeta(ctx context.Context) (code, name string) {
	code = genv.Get("SITE_CODE", "my").String()
	def := "漫隐"
	name = siteconf.Get(ctx, "site_name", def)
	return
}
