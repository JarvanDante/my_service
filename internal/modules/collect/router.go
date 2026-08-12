// Package collect 收藏/点赞模块装配。
package collect

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/collect/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/collect/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台收藏/点赞(需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Operate, ctrl.Delete, ctrl.List)
	})
}
