// Package lottery 抽奖模块装配。
package lottery

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/lottery/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/lottery/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/lottery/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 活动信息公开(挂 AuthOptional: 游客能看奖品与跑马灯, 登录了额外回填
// 我的剩余次数与余额); 抽奖/记录/收货地址需登录, 并额外挂 UserRateLimit 防连点刷奖。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Middleware(middleware.AuthOptional)
		pub.Bind(ctrl.Info)
	})
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Draw, ctrl.My, ctrl.Address, ctrl.MyAddr)
	})
}

// RegisterBackend 活动/奖品 CRUD + 中奖记录 + 收货单发货(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(
		ctrl.ActivityList, ctrl.ActivityCreate, ctrl.ActivityUpdate, ctrl.ActivityDelete,
		ctrl.PrizeList, ctrl.PrizeCreate, ctrl.PrizeUpdate, ctrl.PrizeDelete,
		ctrl.HistoryList, ctrl.AddrList, ctrl.Ship,
	)
}
