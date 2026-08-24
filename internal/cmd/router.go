package cmd

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/dao"
	adminmod "github.com/JarvanDante/my_service/internal/modules/admin"
	aimod "github.com/JarvanDante/my_service/internal/modules/aitask"
	appmod "github.com/JarvanDante/my_service/internal/modules/application"
	checkinmod "github.com/JarvanDante/my_service/internal/modules/checkin"
	collectmod "github.com/JarvanDante/my_service/internal/modules/collect"
	comicsmod "github.com/JarvanDante/my_service/internal/modules/comics"
	commentmod "github.com/JarvanDante/my_service/internal/modules/comment"
	configmod "github.com/JarvanDante/my_service/internal/modules/config"
	couponmod "github.com/JarvanDante/my_service/internal/modules/coupon"
	feedbackmod "github.com/JarvanDante/my_service/internal/modules/feedback"
	groupmod "github.com/JarvanDante/my_service/internal/modules/group"
	finmod "github.com/JarvanDante/my_service/internal/modules/finance"
	lotterymod "github.com/JarvanDante/my_service/internal/modules/lottery"
	mediamod "github.com/JarvanDante/my_service/internal/modules/media"
	msgmod "github.com/JarvanDante/my_service/internal/modules/message"
	novelmod "github.com/JarvanDante/my_service/internal/modules/novel"
	opsmod "github.com/JarvanDante/my_service/internal/modules/ops"
	photomod "github.com/JarvanDante/my_service/internal/modules/photo"
	postmod "github.com/JarvanDante/my_service/internal/modules/post"
	promomod "github.com/JarvanDante/my_service/internal/modules/promo"
	pubmod "github.com/JarvanDante/my_service/internal/modules/publish"
	ranksmod "github.com/JarvanDante/my_service/internal/modules/ranks"
	redeemmod "github.com/JarvanDante/my_service/internal/modules/redeemcode"
	rgmod "github.com/JarvanDante/my_service/internal/modules/redeemgoods"
	searchmod "github.com/JarvanDante/my_service/internal/modules/search"
	statsmod "github.com/JarvanDante/my_service/internal/modules/stats"
	sysmod "github.com/JarvanDante/my_service/internal/modules/system"
	tagmod "github.com/JarvanDante/my_service/internal/modules/tag"
	usermod "github.com/JarvanDante/my_service/internal/modules/user"
	videomod "github.com/JarvanDante/my_service/internal/modules/video"
	walletmod "github.com/JarvanDante/my_service/internal/modules/wallet"
	wdmod "github.com/JarvanDante/my_service/internal/modules/withdrawal"
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
		tagmod.RegisterFront(group)
		redeemmod.RegisterFront(group)
		configmod.RegisterFront(group)
		appmod.RegisterFront(group)
		groupmod.RegisterFront(group)
		collectmod.RegisterFront(group)
		msgmod.RegisterFront(group)
		opsmod.RegisterFront(group)
		commentmod.RegisterFront(group)
		postmod.RegisterFront(group)
		ranksmod.RegisterFront(group)
		rgmod.RegisterFront(group)
		walletmod.RegisterFront(group)
		wdmod.RegisterFront(group)
		couponmod.RegisterFront(group)
		comicsmod.RegisterFront(group)
		novelmod.RegisterFront(group)
		photomod.RegisterFront(group)
		pubmod.RegisterFront(group)
		searchmod.RegisterFront(group)
		lotterymod.RegisterFront(group)
		aimod.RegisterFront(group)
		videomod.RegisterFront(group, dao.NewVideoRepo())
		mediamod.RegisterFront(group, dao.NewMediaRepo())
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
		// 公开: .bnc 解密预览(仅本桶, 给后台 <img>)
		mediamod.RegisterPublic(group, dao.NewMediaRepo())
		// 公开: 支付回调(网关调用, 签名校验在 logic)
		finmod.RegisterCallback(group, dao.NewFinanceRepo())
		// 公开: AI 供应商生成结果回调(验签在 logic, 幂等)
		aimod.RegisterCallback(group)
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
			tagmod.RegisterBackend(perm)
			redeemmod.RegisterBackend(perm)
			configmod.RegisterBackend(perm)
			appmod.RegisterBackend(perm)
			groupmod.RegisterBackend(perm)
			msgmod.RegisterBackend(perm)
			opsmod.RegisterBackend(perm)
			postmod.RegisterBackend(perm)
			commentmod.RegisterBackend(perm)
			ranksmod.RegisterBackend(perm)
			rgmod.RegisterBackend(perm)
			walletmod.RegisterBackend(perm)
			wdmod.RegisterBackend(perm)
			couponmod.RegisterBackend(perm)
			comicsmod.RegisterBackend(perm)
			novelmod.RegisterBackend(perm)
			photomod.RegisterBackend(perm)
			pubmod.RegisterBackend(perm)
			lotterymod.RegisterBackend(perm)
			aimod.RegisterBackend(perm)
			checkinmod.RegisterBackend(perm)
		})
	})
}
