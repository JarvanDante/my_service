// Package message 系统消息模块装配。
package message

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/message/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/message/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/message/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台消息列表/未读/已读(需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.List, ctrl.Unread, ctrl.Read)
	})
}

// RegisterBackend 后台消息发布/管理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
