// Package ranks 排行/热搜模块装配。
package ranks

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/ranks/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/ranks/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/ranks/logic"
)

// RegisterFront 排行/热搜(公开)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.Rank, ctrl.Hot)
}

// RegisterBackend 热搜词管理 + 排行缓存刷新(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.RefreshRank)
}
