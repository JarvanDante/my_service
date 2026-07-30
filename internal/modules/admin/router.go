// Package admin 后台管理员模块装配。
package admin

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/admin/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/admin/domain"
	"github.com/JarvanDante/my_service/internal/modules/admin/logic"
)

// RegisterPublic 公开接口(登录)。
func RegisterPublic(group *ghttp.RouterGroup, repo domain.Repository) {
	group.Bind(backend.New(logic.New(repo)).Login)
}

// RegisterAuthed 仅需登录的接口(退出/查看自身信息), 任何管理员可访问。
func RegisterAuthed(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(ctrl.Logout, ctrl.Info)
}

// RegisterPermManage 需权限校验的接口(角色 / 权限管理)。
func RegisterPermManage(group *ghttp.RouterGroup, repo domain.Repository) {
	ctrl := backend.New(logic.New(repo))
	group.Bind(ctrl.ListRoles, ctrl.ListPerms, ctrl.AddPerm, ctrl.RemovePerm)
}
