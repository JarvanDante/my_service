// Package captcha 开屏图形验证码。
package captcha

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/captcha/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/captcha/logic"
)

// RegisterFront 取码/校验公开(进站前还没有 token)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.Issue, ctrl.Verify)
}
