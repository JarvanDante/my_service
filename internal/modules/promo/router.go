// Package promo 推广/兑换码模块装配(B3)。
package promo

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/promo/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/promo/domain"
	"github.com/JarvanDante/my_service/internal/modules/promo/logic"
)

// RegisterBackend 后台管理接口(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(
		ctrl.CodeList, ctrl.CodeGen, ctrl.CodeVoid, ctrl.CodeLogs,
		// B6 分享 / 拉新
		ctrl.ShareLogs, ctrl.ShareStats,
	)
}
