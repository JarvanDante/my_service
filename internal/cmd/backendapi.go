package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// BackendAPI 后台 API 独立二进制。默认 :8002。
var BackendAPI = gcmd.Command{
	Name:  "backendapi",
	Brief: "后台 API(运营/管理)",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server("backendapi")
		s.SetAddr(cfgAddr(ctx, "backendapi.address", ":8002"))
		mountStatic(ctx, s)
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.BindHandler("/health", middleware.Health)
		mountBackend(s)
		s.Run()
		return nil
	},
}
