package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/rbac"
)

// superAdminCode 超级管理员角色码, 拥有全部权限, 直接放行。
const superAdminCode = "superadmin"

// AdminPerm 后台权限校验(需在 AdminAuth 之后): 超管放行, 其余角色走 Casbin 判定。
func AdminPerm(r *ghttp.Request) {
	adminId := r.GetCtxVar(consts.CtxAdminId).Int64()
	if adminId <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录"))
		return
	}
	// 关联查出角色码
	val, err := g.Model("admin_user au").Ctx(r.GetCtx()).
		LeftJoin("admin_role ar", "ar.id=au.role_id").
		Where("au.id", adminId).
		Fields("ar.code").
		Value()
	if err != nil {
		r.SetError(err)
		return
	}
	roleCode := val.String()
	if roleCode == superAdminCode {
		r.Middleware.Next()
		return
	}
	ok, err := rbac.Enforce(roleCode, r.URL.Path, r.Method)
	if err != nil {
		r.SetError(err)
		return
	}
	if !ok {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "无权限访问该资源"))
		return
	}
	r.Middleware.Next()
}
