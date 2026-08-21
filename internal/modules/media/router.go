// Package media 媒体上传模块(P0: MinIO)。
package media

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/media/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/media/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/media/domain"
	"github.com/JarvanDante/my_service/internal/modules/media/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台用户上传(需登录)。图走整文件加密; 视频走分片续传。
func RegisterFront(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := front.New(logic.New(repo))
	group.Bind(ctrl.Object)
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(
			ctrl.Upload,
			ctrl.MultipartInit,
			ctrl.MultipartPart,
			ctrl.MultipartParts,
			ctrl.MultipartComplete,
			ctrl.MultipartAbort,
		)
	})
}

// RegisterPublic 后台解密预览(给 <img> 用, 仅本桶对象, 无登录)。
func RegisterPublic(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(ctrl.Preview)
}

// RegisterBackend 后台上传接口(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(
		ctrl.Upload,
		ctrl.MultipartInit,
		ctrl.MultipartPresign,
		ctrl.MultipartParts,
		ctrl.MultipartComplete,
		ctrl.MultipartAbort,
	)
}
