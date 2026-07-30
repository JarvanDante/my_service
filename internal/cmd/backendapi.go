package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/dao"
	adminmod "github.com/JarvanDante/my_service/internal/modules/admin"
	usermod "github.com/JarvanDante/my_service/internal/modules/user"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// BackendAPI 后台 API 独立二进制。默认 :8002。
var BackendAPI = gcmd.Command{
	Name:  "backendapi",
	Brief: "后台 API(运营/管理)",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server("backendapi")
		s.SetAddr(cfgAddr(ctx, "backendapi.address", ":8002"))
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.Group("/backend", func(group *ghttp.RouterGroup) {
			// 公开: 管理员登录
			adminmod.RegisterPublic(group, dao.NewAdminRepo())
			// 仅需登录: 退出 / 查看自身信息(任何管理员)
			group.Group("/", func(auth *ghttp.RouterGroup) {
				auth.Middleware(middleware.AdminAuth, middleware.AdminOpLog)
				adminmod.RegisterAuthed(auth, dao.NewAdminRepo())
			})
			// 需权限校验: 角色/权限管理 + 业务接口(超管放行, 其余走 Casbin)
			group.Group("/", func(perm *ghttp.RouterGroup) {
				perm.Middleware(middleware.AdminAuth, middleware.AdminPerm, middleware.AdminOpLog)
				adminmod.RegisterPermManage(perm, dao.NewAdminRepo())
				usermod.RegisterBackend(perm, dao.NewUserRepo())
			})
		})
		s.Run()
		return nil
	},
}
