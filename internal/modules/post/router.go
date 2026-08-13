// Package post 帖子模块装配。
package post

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/post/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/post/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/post/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 列表/详情公开; 发帖/我的/删除需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.List, ctrl.Detail)
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Create, ctrl.My, ctrl.Delete)
	})
}

// RegisterBackend 后台审核/管理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Audit, ctrl.Delete)
}
