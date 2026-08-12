// Package tag 标签模块装配。
package tag

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/tag/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/tag/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/tag/logic"
)

// RegisterFront 前台标签浏览(公开, 无需登录)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.RepoList)
}

// RegisterBackend 后台标签 CRUD(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
