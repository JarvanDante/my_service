// Package redeemcode 兑换码模块装配。
package redeemcode

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/redeemcode/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/redeemcode/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/redeemcode/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台兑换/我的记录(需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Use, ctrl.List)
	})
}

// RegisterBackend 后台建码/管理/使用记录(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Records)
}
