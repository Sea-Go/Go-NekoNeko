package Init

import (
	"context"
	"sea-try-go/service/task/rpc/internal/reward"
	"sea-try-go/service/task/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/proc"
)

func StartInitReward(svc *svc.ServiceContext) {
	ctx, cancel := context.WithCancel(context.Background())

	consumerName := svc.Config.LikeRedis.ConsumerName
	rdb := svc.Rdb
	worker := reward.NewWorker(rdb, svc.PointsClient, consumerName)
	//reclaimer := reward.NewReclaimer(rdb, consumerName)

	go func() {
		_ = worker.Run(ctx)
	}()
	//go reclaimer.Run(ctx)

	proc.AddShutdownListener(func() {
		cancel()
	})
	<-ctx.Done()
}
