// Package withdrawal 提现模块装配。
package withdrawal

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/withdrawal/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/withdrawal/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/withdrawal/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 提现与收款账户全部需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(
			ctrl.Config, ctrl.Apply, ctrl.My, ctrl.Cancel,
			ctrl.CardList, ctrl.CardAdd, ctrl.CardUpdate, ctrl.CardDel,
		)
	})
}

// RegisterBackend 提现审核/打款(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Audit, ctrl.MarkPaid, ctrl.Refund)
}
