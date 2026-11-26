// Package drama 短剧模块装配。
package drama

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/drama/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/drama/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/drama/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 列表/详情/选集/播放公开(带登录态返回解锁态)，单集购买需登录。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New(), logic.NewCategory(), logic.NewModule())
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Middleware(middleware.AuthOptional)
		pub.Bind(ctrl.CategoryList, ctrl.List, ctrl.Detail, ctrl.Episodes, ctrl.ModuleList, ctrl.ModuleRefresh)
	})
	group.Group("/", func(play *ghttp.RouterGroup) {
		play.Middleware(middleware.AuthOptional, middleware.PlayRateLimit)
		play.Bind(ctrl.Play)
	})
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Buy)
	})
}

// RegisterBackend 短剧/剧集/分类/模块/媒资选用(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New(), logic.NewCategory(), logic.NewModule())
	group.Bind(
		ctrl.List, ctrl.Update, ctrl.Delete, ctrl.Audit,
		ctrl.Episodes, ctrl.EpisodeUpdate, ctrl.Sync,
		ctrl.CategoryList, ctrl.CategoryCreate, ctrl.CategoryUpdate, ctrl.CategoryDelete,
		ctrl.ModuleList, ctrl.ModuleCreate, ctrl.ModuleUpdate, ctrl.ModuleDelete,
		ctrl.MediaDramas, ctrl.MediaPick,
	)
}
