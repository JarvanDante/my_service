package cmd

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/dao"
	adminmod "github.com/JarvanDante/my_service/internal/modules/admin"
	checkinmod "github.com/JarvanDante/my_service/internal/modules/checkin"
	feedbackmod "github.com/JarvanDante/my_service/internal/modules/feedback"
	finmod "github.com/JarvanDante/my_service/internal/modules/finance"
	mediamod "github.com/JarvanDante/my_service/internal/modules/media"
	promomod "github.com/JarvanDante/my_service/internal/modules/promo"
	statsmod "github.com/JarvanDante/my_service/internal/modules/stats"
	sysmod "github.com/JarvanDante/my_service/internal/modules/system"
	usermod "github.com/JarvanDante/my_service/internal/modules/user"
	videomod "github.com/JarvanDante/my_service/internal/modules/video"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// mountFront 前台路由统一装配(一体化 Main 与独立 FrontAPI 共用, 单一维护点)。
// 加新前台模块、或大改起新版本(/front/v2), 都只改这里一处, 两个入口自动生效。
func mountFront(s *ghttp.Server) {
	// ---- v1 ----
	s.Group("/front/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.RateLimit)
		usermod.RegisterFront(group, dao.NewUserRepo())
		checkinmod.RegisterFront(group)
		feedbackmod.RegisterFront(group)
	})

	// ---- v2(将来大改时启用; 老的 v1 继续并存)----
	// s.Group("/front/v2", func(group *ghttp.RouterGroup) {
	// 	group.Middleware(middleware.RateLimit)
	// 	usermod.RegisterFrontV2(group, dao.NewUserRepo())
	// })
}

// mountBackend 后台路由统一装配(一体化 Main 与独立 BackendAPI 共用, 单一维护点)。
func mountBackend(s *ghttp.Server) {
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
			videomod.RegisterBackend(perm, dao.NewVideoRepo())
			feedbackmod.RegisterBackend(perm)
		})
	})
}
