// Package search 全站搜索模块装配。
package search

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/search/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/search/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 搜索与联想均公开。挂 AuthOptional 而不是不挂中间件:
// 现在结果里还没有"是否已购"这类个性化字段, 但游客与登录用户走同一入口,
// 先把 userId 认出来放进 ctx, 以后加个性化(已购标记/搜索历史)不用再动路由。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Middleware(middleware.AuthOptional)
		pub.Bind(ctrl.Search, ctrl.Suggest)
	})
}
