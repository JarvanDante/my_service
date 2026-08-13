// Package coupon 优惠券模块装配。
package coupon

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/coupon/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/coupon/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/coupon/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 券模板列表公开; 领券/我的券/可用券需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.Tpls)
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Receive, ctrl.My, ctrl.Available)
	})
}

// RegisterBackend 券模板 CRUD + 发放 + 用户券查询(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Grant, ctrl.Users)
}
