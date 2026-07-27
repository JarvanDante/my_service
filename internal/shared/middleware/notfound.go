package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// NotFound 未匹配到路由时: 若同路径在其他请求方法下存在 -> 提示请求方式不对, 否则接口不存在。
func NotFound(r *ghttp.Request) {
	reqPath := r.URL.Path
	for _, item := range r.Server.GetRoutes() {
		if item.Method == "" || strings.EqualFold(item.Method, "ALL") {
			continue // 跳过中间件/全方法项
		}
		if strings.EqualFold(item.Method, r.Method) {
			continue // 同方法就不是「方法不对」
		}
		if routeMatch(item.Route, reqPath) {
			r.Response.ClearBuffer()
			r.Response.Status = 405
			r.Response.WriteJson(g.Map{
				"code":    405,
				"message": "请求方式不正确, 该接口应使用 " + strings.ToUpper(item.Method),
				"data":    nil,
			})
			return
		}
	}
	r.Response.ClearBuffer()
	r.Response.Status = 404
	r.Response.WriteJson(g.Map{"code": 404, "message": "接口不存在", "data": nil})
}

// routeMatch 简单匹配 GoFrame 路由模式({name}/:name/*any) 与请求路径。
func routeMatch(pattern, path string) bool {
	ps := strings.Split(strings.Trim(pattern, "/"), "/")
	xs := strings.Split(strings.Trim(path, "/"), "/")
	for i, seg := range ps {
		if strings.HasPrefix(seg, "*") {
			return true // 通配剩余
		}
		if i >= len(xs) {
			return false
		}
		if strings.HasPrefix(seg, "{") || strings.HasPrefix(seg, ":") {
			continue // 参数段任意匹配
		}
		if seg != xs[i] {
			return false
		}
	}
	return len(ps) == len(xs)
}
