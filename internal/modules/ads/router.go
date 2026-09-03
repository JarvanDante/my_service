package ads

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/ads/controller/front"
)

func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New()
	group.Bind(ctrl.List, ctrl.Event)
}
