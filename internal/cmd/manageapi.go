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

// ManageAPI 总后台 API 独立二进制。默认 :8003。
var ManageAPI = gcmd.Command{
	Name:  "manageapi",
	Brief: "总后台 API(平台超管)",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server("manageapi")
		s.SetAddr(cfgAddr(ctx, "manageapi.address", ":8003"))
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.Group("/manage", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth)
			usermod.RegisterManage(group, dao.NewUserRepo())
		})
		s.Run()
		return nil
	},
}
