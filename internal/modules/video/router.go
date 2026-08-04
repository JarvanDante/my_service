// Package video 视频管理模块(P1)。
package video

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/video/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/video/domain"
	"github.com/JarvanDante/my_service/internal/modules/video/logic"
)

// RegisterBackend 后台视频 CRUD(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Status)
}
