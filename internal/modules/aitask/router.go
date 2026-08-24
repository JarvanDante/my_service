// Package aitask AI 生成任务模块装配(换脸/脱衣/文生图/图生视频/文生小说/AI对话共用)。
package aitask

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/JarvanDante/my_service/internal/modules/aitask/controller/backend"
	"github.com/JarvanDante/my_service/internal/modules/aitask/controller/front"
	"github.com/JarvanDante/my_service/internal/modules/aitask/logic"
	"github.com/JarvanDante/my_service/internal/shared/middleware"
)

// RegisterFront 前台: 模板列表公开(未登录也要能看到玩法与价格);
// 提交/查询/取消都涉及金币与个人数据, 必须登录, 并挂用户级限流
// (AI 生成是重成本操作, 不限流会被脚本刷爆供应商配额和用户余额)。
func RegisterFront(group *ghttp.RouterGroup) {
	ctrl := front.New(logic.New())
	group.Group("/", func(pub *ghttp.RouterGroup) {
		pub.Bind(ctrl.Templates)
	})
	group.Group("/", func(auth *ghttp.RouterGroup) {
		auth.Middleware(middleware.Auth, middleware.UserRateLimit)
		auth.Bind(ctrl.Submit, ctrl.Task, ctrl.Tasks, ctrl.Cancel)
	})
}

// RegisterBackend 后台: 模板 CRUD + 任务列表 + 重新提交 + 人工退款(挂权限组)。
func RegisterBackend(group *ghttp.RouterGroup) {
	ctrl := backend.New(logic.New())
	group.Bind(
		ctrl.TemplateList, ctrl.TemplateCreate, ctrl.TemplateUpdate, ctrl.TemplateDelete,
		ctrl.TaskList, ctrl.TaskRetry, ctrl.TaskRefund, ctrl.TaskDelete,
	)
}

// RegisterCallback 供应商回调(公开, 必须挂在**不鉴权**的分组上)。
//
// 为什么单独一个注册函数: 回调是第三方服务器发起的, 它没有、也不可能有我们的后台管理员
// token —— 挂进 AdminAuth 分组只会被 401 挡在门外, 然后供应商无限重推。
// 它的身份证明是请求里的签名(logic 内校验), 与 finance 的支付回调是同一套做法。
func RegisterCallback(group *ghttp.RouterGroup) {
	group.Bind(backend.New(logic.New()).Callback)
}
