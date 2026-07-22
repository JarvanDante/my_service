// Package user 用户模块装配。按门面提供注册函数, 共享同一份 logic/domain。
package user

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/user/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/user/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/user/controller/manage"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/logic"
)

// RegisterFront 前台路由。
func RegisterFront(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(front.New(logic.New(repo)))
}

// RegisterBackend 后台路由。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(backend.New(logic.New(repo)))
}

// RegisterManage 总后台路由。
func RegisterManage(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(manage.New(logic.New(repo)))
}
