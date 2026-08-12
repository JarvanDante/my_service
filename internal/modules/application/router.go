// Package application 推广应用模块装配。
package application

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/application/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/application/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/application/logic"
)

// RegisterFront 前台应用列表/点击上报(公开)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.List, ctrl.Click)
}

// RegisterBackend 后台应用管理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
