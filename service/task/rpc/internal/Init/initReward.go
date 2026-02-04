package Init

import (
	"context"
	"os"
	"sea-try-go/service/task/rpc/internal/reward"
	"sea-try-go/service/task/rpc/internal/svc"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/proc"
)

func StartInitReward(svc *svc.ServiceContext) {
	ctx, cancel := context.WithCancel(context.Background())

	consumerName := svc.Config.LikeRedis.ConsumerName
	rdb := svc.Rdb
	worker := reward.NewWorker(rdb, pointsClient, consumerName)
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
