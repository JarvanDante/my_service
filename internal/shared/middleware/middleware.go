// Package middleware 全局/分组中间件。
package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// Auth 鉴权: 解析 Authorization 里的 token, 成功则把 userId 写入 ctx, 失败返回未授权。
func Auth(r *ghttp.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.Get("token").String() // 兼容用 query/form 传 token
	}
	uid, err := kit.ParseToken(r.GetCtx(), token)
	if err != nil || uid <= 0 {
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "未登录或登录已失效"))
		return
	}
	r.SetCtxVar(consts.CtxUserId, uid)
	r.Middleware.Next()
}

// RateLimit 前台限流占位。
func RateLimit(r *ghttp.Request) {
	// TODO: 基于 Redis 的令牌桶 / 滑动窗口
	r.Middleware.Next()
}
