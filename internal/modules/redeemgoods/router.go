// Package redeemgoods 商品兑换模块装配。
package redeemgoods

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 商品列表公开; 兑换/历史需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.List)
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Exchange, ctrl.History)
	})
}

// RegisterBackend 商品管理 + 兑换记录(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Orders)
}
