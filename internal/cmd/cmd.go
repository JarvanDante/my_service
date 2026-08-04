// Package cmd 各二进制的启动装配。业务依赖(dao)在此注入各模块。
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/JarvanDante/my_service/internal/dao"
	"github.com/JarvanDante/my_service/internal/shared/middleware"

	adminmod "github.com/JarvanDante/my_service/internal/modules/admin"
	finmod "github.com/JarvanDante/my_service/internal/modules/finance"
	mediamod "github.com/JarvanDante/my_service/internal/modules/media"
	promomod "github.com/JarvanDante/my_service/internal/modules/promo"
	statsmod "github.com/JarvanDante/my_service/internal/modules/stats"
	sysmod "github.com/JarvanDante/my_service/internal/modules/system"
	usermod "github.com/JarvanDante/my_service/internal/modules/user"
)

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
		s.Use(middleware.CORS, ghttp.MiddlewareHandlerResponse)
		s.BindStatusHandler(404, middleware.NotFound)
		s.BindHandler("/health", middleware.Health)

		s.Group("/front", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.RateLimit)
			usermod.RegisterFront(group, dao.NewUserRepo())
		})
		s.Group("/backend", func(group *ghttp.RouterGroup) {
			// 公开: 管理员登录
			adminmod.RegisterPublic(group, dao.NewAdminRepo())
			// 公开: 支付回调(网关调用, 签名校验在 logic)
			finmod.RegisterCallback(group, dao.NewFinanceRepo())
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
				finmod.RegisterBackend(perm, dao.NewFinanceRepo())
				promomod.RegisterBackend(perm, dao.NewPromoRepo())
				sysmod.RegisterBackend(perm, dao.NewSystemRepo())
				statsmod.RegisterBackend(perm, dao.NewStatsRepo())
				mediamod.RegisterBackend(perm, dao.NewMediaRepo())
			})
		})
		s.Group("/manage", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth)
			usermod.RegisterManage(group, dao.NewUserRepo())
		})

		registerCronJobs(ctx)
		s.Run()
		return nil
	},
}
