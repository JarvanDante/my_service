package channel

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/channel/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/channel/logic"
)

func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
