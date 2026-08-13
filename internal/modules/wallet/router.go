// Package wallet 钱包模块装配。
package wallet

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/wallet/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/wallet/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/wallet/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 余额/流水均需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Balance, ctrl.Waters)
	})
}

// RegisterBackend 全站流水 + 人工调账(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.Logs, ctrl.Adjust)
}
