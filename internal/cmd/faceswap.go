package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/modules/aitask/logic"
	"github.com/JarvanDante/my_service/internal/mq"
	"github.com/JarvanDante/my_service/internal/shared/aiprotocol"
	"github.com/JarvanDante/my_service/internal/shared/aiprovider"
)

// startFaceSwapWorker 装配本地 AI 工人(换脸/去衣): Kafka 投递 + 结果消费。
func startFaceSwapWorker(ctx context.Context) {
	bus := mq.Init(ctx)
	aiprovider.SetFaceSwapPublisher(bus.PublishJob)
	if !bus.Enabled() {
		g.Log().Warning(ctx, "kafka disabled, skip faceswap result consumer")
		return
	}
	svc := logic.New()
	go func() {
		err := bus.ConsumeResults(ctx, func(ctx context.Context, res aiprotocol.ResultMessage) error {
			return svc.HandleWorkerResult(ctx, res.JobID, res.Status, res.OutputURL, res.OutputKey, res.Error)
		})
		if err != nil && ctx.Err() == nil {
			g.Log().Errorf(ctx, "faceswap result consumer stopped: %v", err)
		}
	}()
}
