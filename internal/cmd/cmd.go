// Package cmd 各二进制的启动装配。业务依赖(dao)在此注入各模块。
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// mountStatic 挂内置静态资源(默认头像包等)。目录不存在时跳过并告警,
// 不让非关键的静态资源拦住服务启动(比如镜像里漏拷了 resource/public)。
func mountStatic(ctx context.Context, s *ghttp.Server) {
	const dir = "resource/public"
	if !gfile.Exists(dir) {
		g.Log().Warningf(ctx, "静态资源目录 %s 不存在, 跳过 /static 挂载(默认头像将 404)", dir)
		return
	}
	s.AddStaticPath("/static", dir)
}

// cfgAddr 读监听地址, 缺省回退默认。
func cfgAddr(ctx context.Context, key, def string) string {
	if v, err := g.Cfg().Get(ctx, key); err == nil && !v.IsNil() && v.String() != "" {
		return v.String()
	}
	return def
}

// Main 本地开发一体化入口: 单进程挂载全部门面 + 定时任务, 方便 gf run。
// 生产环境请用 app/ 下的独立二进制以实现进程隔离。
var Main = gcmd.Command{
	Name:  "main",
	Brief: "漫隐 API · 一体化开发入口",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server()
		mountStatic(ctx, s)
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.BindHandler("/health", middleware.Health)

		mountFront(s)
		mountBackend(s)

		registerCronJobs(ctx)
		s.Run()
		return nil
	},
}
