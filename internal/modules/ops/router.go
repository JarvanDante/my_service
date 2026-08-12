// Package ops 运营配置模块装配(公告/跳转位/敏感词)。
package ops

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/ops/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/ops/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/ops/logic"
)

// RegisterFront 前台公告/跳转位(公开)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.Announcement, ctrl.Jumptab)
}

// RegisterBackend 后台运营配置管理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(
		ctrl.AnnList, ctrl.AnnCreate, ctrl.AnnUpdate, ctrl.AnnDelete,
		ctrl.JtList, ctrl.JtCreate, ctrl.JtUpdate, ctrl.JtDelete,
		ctrl.FwList, ctrl.FwAdd, ctrl.FwDelete,
	)
}
