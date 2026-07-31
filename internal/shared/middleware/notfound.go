package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// NotFound 未匹配到路由时的统一处理:
//  1. 当前方法下存在「比请求路径多一个参数段」的路由 -> 少传了路径参数(如 DELETE /admins 少了 {id})
//  2. 同路径在其他请求方法下存在 -> 请求方式不对
//  3. 否则 -> 接口不存在
func NotFound(r *ghttp.Request) {
	reqPath := r.URL.Path
	routes := r.Server.GetRoutes()

	// 1) 少传路径参数: 当前方法下有 reqPath + "/{param}" 形态的路由
	for _, item := range routes {
		if item.Method == "" || strings.EqualFold(item.Method, "ALL") {
			continue
		}
		if !strings.EqualFold(item.Method, r.Method) {
			continue
		}
		if isParamChild(item.Route, reqPath) {
			r.Response.ClearBuffer()
			r.Response.Status = 400
			r.Response.WriteJson(g.Map{
				"code":    400,
				"message": "缺少路径参数, 请在 URL 中补全 (如 " + item.Route + ")",
				"data":    nil,
			})
			return
		}
	}

	// 2) 方法不对: 同路径存在其他方法
	for _, item := range routes {
		if item.Method == "" || strings.EqualFold(item.Method, "ALL") {
			continue
		}
		if strings.EqualFold(item.Method, r.Method) {
			continue
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

	// 3) 接口不存在
	r.Response.ClearBuffer()
	r.Response.Status = 404
	r.Response.WriteJson(g.Map{"code": 404, "message": "接口不存在", "data": nil})
}

// isParamChild 判断 pattern 是否为 path 后面「恰好多一个参数段」的路由。
// 例: pattern=/backend/admins/{id}, path=/backend/admins -> true。
func isParamChild(pattern, path string) bool {
	ps := strings.Split(strings.Trim(pattern, "/"), "/")
	xs := strings.Split(strings.Trim(path, "/"), "/")
	if len(ps) != len(xs)+1 {
		return false
	}
	last := ps[len(ps)-1]
	if !strings.HasPrefix(last, "{") && !strings.HasPrefix(last, ":") {
		return false // 多出的那段必须是参数段
	}
	for i, seg := range xs {
		p := ps[i]
		if strings.HasPrefix(p, "{") || strings.HasPrefix(p, ":") {
			continue
		}
		if p != seg {
			return false
		}
	}
	return true
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
