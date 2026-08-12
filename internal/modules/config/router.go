// Package config 基础配置模块装配。
package config

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/config/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/config/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/config/logic"
)

// RegisterFront 前台配置读取(公开, 客户端启动第一个协议)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Bind(ctrl.Info, ctrl.Check)
}

// RegisterBackend 后台配置 KV 管理(挂在权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(ctrl.List, ctrl.Create, ctrl.Update, ctrl.Delete)
}
