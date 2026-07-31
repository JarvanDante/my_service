// Package finance 财务模块装配(B2)。
package finance

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/finance/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/finance/domain"
	"github.com/JarvanDante/my_service/internal/modules/finance/logic"
)

// RegisterBackend 后台管理接口(挂权限组): 套餐 CRUD + 订单 + 流水。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(
		ctrl.RechargePkgList, ctrl.RechargePkgCreate, ctrl.RechargePkgUpdate, ctrl.RechargePkgDelete,
		ctrl.VipPkgList, ctrl.VipPkgCreate, ctrl.VipPkgUpdate, ctrl.VipPkgDelete,
		ctrl.Orders, ctrl.BalanceLogs,
	)
}

// RegisterCallback 支付回调(公开, 由支付网关调用, 签名校验在 logic)。
func RegisterCallback(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(backend.New(logic.New(repo)).PayCallback)
}
