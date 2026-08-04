// Package media 媒体上传模块(P0: MinIO)。
package media

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/media/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/media/domain"
	"github.com/JarvanDante/my_service/internal/modules/media/logic"
)

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
