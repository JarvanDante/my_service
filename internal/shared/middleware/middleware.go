// Package middleware 全局中间件。
package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// CORS 跨域。
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// Response 统一响应处理(可按需扩展成标准返回结构)。
func Response(r *ghttp.Request) {
	r.Middleware.Next()
}

// Auth 鉴权占位: 后续接入 JWT 校验 + Casbin 权限。
func Auth(r *ghttp.Request) {
	// TODO: 解析 JWT -> 校验 Casbin 策略
	r.Middleware.Next()
}
