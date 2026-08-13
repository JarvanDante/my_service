// Package comics 漫画模块装配。
package comics

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/comics/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/comics/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/comics/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 列表/详情/目录/阅读公开(带 token 则识别已购与解锁态); 购买需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Middleware(middleware.AuthOptional)
		pub.Bind(ctrl.List, ctrl.Detail, ctrl.Chapters, ctrl.Read, ctrl.MayLike)
	})
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Buy)
	})
}

// RegisterBackend 作品 CRUD + 上下架 + 章节 CRUD(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(
		ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Audit,
		ctrl.Chapters, ctrl.ChapterCreate, ctrl.ChapterUpdate, ctrl.ChapterDelete,
	)
}
