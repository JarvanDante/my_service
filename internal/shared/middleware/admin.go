package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

// AdminAuth 后台鉴权: 解析管理员 token, 成功则把 adminId 写入 ctx。
func AdminAuth(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.Get("token").String()
	}
	adminId, err := kit.ParseAdminToken(r.GetCtx(), token)
	if err != nil || adminId <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录或登录已失效"))
		return
	}
	r.SetCtxVar(consts.CtxAdminId, adminId)
	r.Middleware.Next()
}

// AdminOpLog 记录后台写操作日志(POST/PUT/DELETE)。
func AdminOpLog(r *ghttp.Request) {
	r.Middleware.Next()
	m := r.Method
	if m == "POST" || m == "PUT" || m == "DELETE" {
		adminId := r.GetCtxVar(consts.CtxAdminId).Int64()
		if adminId > 0 {
			_, _ = g.Model("admin_log").Ctx(r.GetCtx()).Data(g.Map{
				"admin_id": adminId,
				"method":   m,
				"path":     r.URL.Path,
				"ip":       r.GetClientIp(),
			}).Insert()
		}
	}
}
