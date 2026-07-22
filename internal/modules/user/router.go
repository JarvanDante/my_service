// Package user 用户模块装配。按门面提供注册函数, 共享同一份 logic/domain。
package user

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/user/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/user/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/user/controller/manage"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台路由: 公开(login) + 需登录(其余)。
func RegisterFront(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := front.New(logic.New(repo))
	// 公开接口(无需 token)
	group.Bind(ctrl.Login)
	// 需登录接口
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth)
		auth.Bind(ctrl.Info)
	})
}

// RegisterBackend 后台路由(整体需鉴权, 由入口挂 Auth)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(backend.New(logic.New(repo)))
}

// RegisterManage 总后台路由。
func RegisterManage(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(manage.New(logic.New(repo)))
}
