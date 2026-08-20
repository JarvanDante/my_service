// Package video 视频管理模块(P1)。
package video

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/video/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/video/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/video/domain"
	"github.com/JarvanDante/my_service/internal/modules/video/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台视频浏览(公开)。repo 由调用方注入, 与 RegisterBackend 同一套仓储,
// 前后台共用一个 logic 实例语义(各自 New 一份, 无状态)。
// 挂 AuthOptional: 现在列表/详情不区分登录态, 但视频后续要接付费墙(paywall.MediaVideo),
// 届时需要 ctx 里的 userId 才能标记已购, 先把中间件位置留好。
func RegisterFront(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := front.New(logic.New(repo), logic.NewCategory(), logic.NewCategoryTable("cartoon_category"))
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Middleware(middleware.AuthOptional)
		pub.Bind(ctrl.List, ctrl.Detail, ctrl.CategoryList, ctrl.CartoonList, ctrl.CartoonCategoryList)
	})
}

// RegisterBackend 后台视频 CRUD + 分类 CRUD(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo), logic.NewCategory(), logic.NewCategoryTable("cartoon_category"))
	group.Bind(
		ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete, ctrl.Status,
		ctrl.MediaAssets, ctrl.MediaPick, ctrl.SyncMedia,
		ctrl.CategoryList, ctrl.CategoryCreate, ctrl.CategoryUpdate, ctrl.CategoryDelete,
		ctrl.CartoonCategoryList, ctrl.CartoonCategoryCreate, ctrl.CartoonCategoryUpdate, ctrl.CartoonCategoryDelete,
	)
}
