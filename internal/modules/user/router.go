package user

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/user/controller"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/logic"
)

// Register 由 cmd 调用: 注入仓储 -> 装配本模块 -> 绑定路由。
func Register(group *ghttp.RouterGroup, repo domain.Repository) {
	svc := logic.New(repo)
	ctrl := controller.New(svc)
	group.Bind(ctrl)
}
