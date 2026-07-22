// Package cmd 应用启动与依赖装配点。唯一把 dao 实现注入模块、并注册路由的地方。
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/dao"
	"github.com/JarvanDante/my_service/internal/shared/middleware"

	usermod "github.com/JarvanDante/my_service/internal/modules/user"
)

var Main = gcmd.Command{
	Name:  "main",
	Brief: "漫隐 API 服务",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		s := g.Server()

		// 全局中间件
		s.Use(middleware.CORS, middleware.Response)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth) // 需鉴权的分组

			// ==== 逐模块装配: 注入仓储 -> 注册路由 ====
			usermod.Register(group, dao.NewUserRepo())
			// ordermod.Register(group, dao.NewOrderRepo())  // 新增模块照此追加一行
		})

		s.Run()
		return nil
	},
}
