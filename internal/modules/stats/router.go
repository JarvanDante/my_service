// Package stats 统计模块装配(B8)。
package stats

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/stats/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/stats/domain"
	"github.com/JarvanDante/my_service/internal/modules/stats/logic"
)

// RegisterBackend 后台统计接口(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(ctrl.Overview, ctrl.UserTrend, ctrl.RechargeTrend, ctrl.ChannelStats)
}
