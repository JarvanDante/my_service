package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// FrontAPI 前台 API 独立二进制。默认 :8001。
var FrontAPI = gcmd.Command{
	Name:  "frontapi",
	Brief: "前台 API(面向 C 端)",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server("frontapi")
		s.SetAddr(cfgAddr(ctx, "frontapi.address", ":8001"))
		s.AddStaticPath("/static", "resource/public") // 内置静态资源(默认头像包等)
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.BindHandler("/health", middleware.Health)
		mountFront(s)
		s.Run()
		return nil
	},
}
