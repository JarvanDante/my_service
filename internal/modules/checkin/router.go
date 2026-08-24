// Package checkin 签到模块装配。
package checkin

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/checkin/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/checkin/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/checkin/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台签到(需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Click, ctrl.Info)
	})
}

// RegisterBackend 签到配置(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.GetConfig, ctrl.SaveConfig, ctrl.RewardList, ctrl.SaveRewards)
}
