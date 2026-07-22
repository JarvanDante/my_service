package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gcron"

	"github.com/JarvanDante/my_service/internal/dao"
	"github.com/JarvanDante/my_service/internal/modules/user/logic"
)

// Cron 定时任务独立二进制。无 HTTP 端口, 与 API 进程隔离。
var Cron = gcmd.Command{
	Name:  "cron",
	Brief: "定时任务",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		registerCronJobs(ctx)
		g.Log().Info(ctx, "cron started")
		select {} // 阻塞常驻; 生产可接信号做优雅退出
	},
}

// registerCronJobs 集中注册定时任务(cron 二进制与一体化 dev 入口共用)。
func registerCronJobs(ctx context.Context) {
	userLogic := logic.New(dao.NewUserRepo())
	_ = userLogic

	// GoFrame cron 为 6 段(含秒)。示例: 每天 03:00 执行。
	_, _ = gcron.Add(ctx, "0 0 3 * * *", func(ctx context.Context) {
		g.Log().Info(ctx, "daily job running...")
		// TODO: 调用 userLogic / 其他模块 service 完成批处理
	}, "daily-user-job")
}
