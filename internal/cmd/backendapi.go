package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/dao"
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
			group.Middleware(middleware.Auth)
			usermod.RegisterBackend(group, dao.NewUserRepo())
		})
		s.Run()
		return nil
	},
}
