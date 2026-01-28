package Init

import (
	"context"
	"log"
	"sea-try-go/service/task/rpc/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
)

type PrintConsumer struct {
}

func (c PrintConsumer) Consume(ctx context.Context, key string, val string) error {
	log.Printf("got msg key=%s val=%s\n", key, val)
	return nil
}

func StartTaskKafkaConsumer(svcCtx *svc.ServiceContext) {
	ctx, cancel := context.WithCancel(context.Background())

	brokers := svcCtx.Config.Kafka.Brokers
	topic := svcCtx.Config.Kafka.Topic
	group := svcCtx.Config.Kafka.Group

	q := kq.MustNewQueue(kq.KqConf{
		Brokers:    brokers,
		Topic:      topic,
		Group:      group,
		Consumers:  1,
		Offset:     "latest",
		Processors: 1,
	}, PrintConsumer{})

	// 退出时停掉 consumer
	proc.AddShutdownListener(func() {
		log.Println("kafka consumer stopping...")
		cancel()
		q.Stop()
	})

	q.Start()

	<-ctx.Done()
}
