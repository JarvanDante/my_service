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

// FrontAPI 前台 API 独立二进制。默认 :8001。
var FrontAPI = gcmd.Command{
	Name:  "frontapi",
	Brief: "前台 API(面向 C 端)",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server("frontapi")
		s.SetAddr(cfgAddr(ctx, "frontapi.address", ":8001"))
		s.Use(middleware.CORS, middleware.Response)
		s.Group("/front", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.RateLimit)
			usermod.RegisterFront(group, dao.NewUserRepo())
		})
		s.Run()
		return nil
	},
}
