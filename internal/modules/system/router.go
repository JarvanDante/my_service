// Package system 系统模块装配(B7)。
package system

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/system/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/system/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/system/domain"
	"github.com/JarvanDante/my_service/internal/modules/system/logic"
)

func RegisterFront(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := front.New(logic.New(repo))
	group.Bind(ctrl.List)
}

// RegisterBackend 后台管理接口(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(
		ctrl.Push, ctrl.NoticeList, ctrl.NoticeStatus,
		ctrl.CustomerUrlGet, ctrl.CustomerUrlPut,
	)
}
