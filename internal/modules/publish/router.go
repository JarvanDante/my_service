// Package publish UGC投稿模块装配。
package publish

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/publish/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/publish/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/publish/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 投稿全部是用户私域操作, 一律需要登录并限流(防脚本刷投稿)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Submit, ctrl.My, ctrl.Cancel)
	})
}

// RegisterBackend 投稿列表 + 审核(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Audit)
}
