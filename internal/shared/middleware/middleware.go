// Package middleware 全局/分组中间件。
package middleware

import (
	"net/http"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/genv"

	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
	"github.com/JarvanDante/my_service/internal/shared/ratelimit"
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

// RateLimit 前台按客户端 IP 限流(登录前即生效, 防暴力登录/刷接口)。
func RateLimit(r *ghttp.Request) {
	ctx := r.Context()
	ip := r.GetClientIp()
	perMin := g.Cfg().MustGet(ctx, "front_ratelimit.ip_per_min", 300).Int()
	perHour := g.Cfg().MustGet(ctx, "front_ratelimit.ip_per_hour", 5000).Int()
	if ok, retry, reason := ratelimit.Allow(ctx, "ip", ip, perMin, perHour); !ok {
		g.Log().Warningf(ctx, "前台IP限流命中 ip=%s path=%s: %s", ip, r.URL.Path, reason)
		rlDeny(r, retry, reason)
		return
	}
	r.Middleware.Next()
}

// UserRateLimit 按登录用户限流(须在 Auth 之后挂载, 依赖 ctx 里的 userId)。
// 用于压制"注册账号后批量刷接口/刷播放地址"的滥用。
func UserRateLimit(r *ghttp.Request) {
	ctx := r.Context()
	uid := r.GetCtxVar(consts.CtxUserId).Int64()
	if uid <= 0 {
		r.Middleware.Next()
		return
	}
	perMin := g.Cfg().MustGet(ctx, "front_ratelimit.user_per_min", 120).Int()
	perHour := g.Cfg().MustGet(ctx, "front_ratelimit.user_per_hour", 1500).Int()
	if ok, retry, reason := ratelimit.Allow(ctx, "user", strconv.FormatInt(uid, 10), perMin, perHour); !ok {
		g.Log().Warningf(ctx, "前台用户限流命中 uid=%d ip=%s path=%s: %s", uid, r.GetClientIp(), r.URL.Path, reason)
		rlDeny(r, retry, reason)
		return
	}
	r.Middleware.Next()
}

// rlDeny 统一返回 429 + Retry-After。
func rlDeny(r *ghttp.Request, retryAfter int, reason string) {
	r.Response.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	r.Response.WriteStatus(http.StatusTooManyRequests)
	r.Response.WriteJsonExit(g.Map{"code": 429, "message": reason, "data": nil})
}

// Health 公开健康检查端点(网关/LB/总后台探活)。
func Health(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"status": "ok",
		"site":   genv.Get("SITE_CODE", "my").String(),
	})
}
