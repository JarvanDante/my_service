package banner

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/banner/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/banner/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/banner/logic"
)

func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.List)
}

func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
