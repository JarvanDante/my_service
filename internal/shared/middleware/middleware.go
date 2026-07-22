// Package middleware 全局中间件。
package middleware

import "github.com/gogf/gf/v2/net/ghttp"

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Response(r *ghttp.Request) {
	r.Middleware.Next()
}

// Auth 鉴权占位: 后续 JWT + Casbin。后台/总后台专用。
func Auth(r *ghttp.Request) {
	// TODO: 解析 JWT -> 校验 Casbin 策略
	r.Middleware.Next()
}

// RateLimit 前台限流占位。
func RateLimit(r *ghttp.Request) {
	// TODO: 基于 Redis 的令牌桶 / 滑动窗口
	r.Middleware.Next()
}
