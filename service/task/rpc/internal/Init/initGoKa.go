package Init

import (
	"context"
	"log"
	"sea-try-go/service/task/rpc/internal/svc"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/zeromicro/go-zero/core/proc"
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
	brokerTask   []string
	inTopicTask  string
	outTopicTask string
	groupTask    string
	topicDLQ     = "raw-logs.dlq"
	countTable   = "like_counts"
)

func process(ctx goka.Context, msg any) { //核心处理逻辑
	_ = msg.([]byte)

	var cur int64
	if v := ctx.Value(); v != nil {
		cur = v.(int64)
	}
	cur++
	ctx.SetValue(cur)
	ctx.Emit(goka.Stream(outTopicTask), ctx.Key(), cur)
}

func StartTaskGoKa(svcCtx *svc.ServiceContext) {

	ctx, cancel := context.WithCancel(context.Background())

	brokerTask = svcCtx.Config.Kafka.Brokers
	groupTask = svcCtx.Config.Kafka.GroupGoKa
	inTopicTask = svcCtx.Config.Kafka.InTopic
	outTopicTask = svcCtx.Config.Kafka.OutTopic

	g := goka.DefineGroup(
		goka.Group(groupTask),
		goka.Input(goka.Stream(inTopicTask), new(codec.Bytes), process),
		goka.Persist(new(codec.Int64)),
		goka.Output(goka.Stream(outTopicTask), new(codec.Int64)),
	)

	cfg := goka.DefaultConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	//修复topic复制因子rfactor=2的问题
	tmCfg := goka.NewTopicManagerConfig()
	tmCfg.Stream.Replication = 1
	tmCfg.Table.Replication = 1

	p, err := goka.NewProcessor(brokerTask, g,
		goka.WithConsumerGroupBuilder(goka.ConsumerGroupBuilderWithConfig(cfg)),
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithConfig(cfg, tmCfg)),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("like aggregator started: in=%s table=%s out=%s\n", inTopicTask, countTable, outTopicTask)
	if err := p.Run(ctx); err != nil {
		log.Fatal(err)
	}

	proc.AddShutdownListener(func() {
		cancel()
	})
	<-ctx.Done()
}
