// Package feedback 意见反馈模块装配。
package feedback

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/feedback/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/feedback/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/feedback/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台提交(需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Add)
	})
}

// RegisterBackend 后台列表/处理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Handle)
}
