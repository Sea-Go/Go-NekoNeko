package Init

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sea-try-go/service/task/rpc/internal/svc"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

/*var (
	brokers  = []string{}
	group    string
	inTopic  string
	outTopic string
)*/

// 日志结构体
type LogEvent struct {
	TS      string `json:"ts"`      // 时间戳 2026-01-22T10:11:12.123Z
	Service string `json:"service"` // 服务名 like-service/video-service
	//Position string `json:"position"` // 位置
	EventID   string `json:"event_id"`   // 全局唯一事件ID，用于幂等
	UserID    string `json:"user_id"`    // 用户ID
	TraceID   string `json:"trace_id"`   // 链路追踪ID
	ArticleID string `json:"article_id"` //文章ID
	Msg       string `json:"msg"`        // 信息
}

type DLQEvent struct {
	RawLog     []byte `json:"raw_log"`
	FailReason string `json:"fail_reason"`
	ProcessTS  string `json:"process_ts"`
}

var (
	broker     []string
	inTopic    string
	outTopic   string
	group      string
	topicDLQ   = "raw-logs.dlq"
	countTable = "like_counts"
)

/*type Int64Codec struct{}

func (c *Int64Codec) Encode(value interface{}) ([]byte, error) {
	v := value.(int64)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b, nil
}

func (c *Int64Codec) Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return int64(0), nil
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}
*/
func process(ctx goka.Context, msg any) {
	_ = msg.([]byte)

	var cur int64
	if v := ctx.Value(); v != nil {
		cur = v.(int64)
	}
	cur++

	ctx.SetValue(cur)

	ctx.Emit(goka.Stream(outTopic), ctx.Key(), cur)
}

func startTaskGoKa(svcCtx *svc.ServiceContext) {

	ctx, cancel := context.WithCancel(context.Background())

	broker = svcCtx.Config.Kafka.Brokers
	group = svcCtx.Config.Kafka.Group
	inTopic = svcCtx.Config.Kafka.InTopic
	outTopic = svcCtx.Config.Kafka.OutTopic

	g := goka.DefineGroup(
		goka.Group(group),
		goka.Input(goka.Stream(inTopic), new(codec.Bytes), process),
		goka.Persist(new(codec.Int64)),
		goka.Output(goka.Stream(outTopic), new(codec.String)),
	)

	cfg := goka.DefaultConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	p, err := goka.NewProcessor(broker, g,
		goka.WithConsumerGroupBuilder(goka.ConsumerGroupBuilderWithConfig(cfg)),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	go func() {
		ch := make(chan os.Signal, 2)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()
	log.Printf("like aggregator started: in=%s table=%s out=%s\n", inTopic, countTable, outTopic)
	if err := p.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
