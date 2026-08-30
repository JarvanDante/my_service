package kingkong

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/kingkong/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/kingkong/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/kingkong/logic"
)

func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.List)
}

func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
