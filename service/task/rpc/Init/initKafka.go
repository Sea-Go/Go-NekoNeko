package Init

import (
	"context"
	"fmt"
	"log"
	"sea-try-go/service/task/rpc/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
)

type PrintConsumer struct {
}

func (c PrintConsumer) Consume(ctx context.Context, key string, val string) error {
	fmt.Printf("got msg key=%s val=%s\n", key, val)
	return nil
}

func StartTaskKafkaConsumer(svcCtx *svc.ServiceContext) {
	ctx, cancel := context.WithCancel(context.Background())
	proc.AddShutdownListener(func() {
		log.Println("start kafka consumer")
		cancel()
	})

	brokers := svcCtx.Config.Kafka.Brokers
	topic := svcCtx.Config.Kafka.Topic
	group := svcCtx.Config.Kafka.Group

	q := kq.MustNewQueue(kq.KqConf{
		Brokers:   brokers,
		Topic:     topic,
		Group:     group,
		Consumers: 1,
		Offset:    "latest",
	}, PrintConsumer{})

	q.Start()
	defer q.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("start kafka consumer")
			return
		default:

		}
	}
}
