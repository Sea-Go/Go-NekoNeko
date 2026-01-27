package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/service/points/rpc/internal/svc"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
)

type KafkaConsumer struct {
	sCtx *svc.ServiceContext
}

func NewKafkaConsumer(sCtx *svc.ServiceContext) *KafkaConsumer {
	kc := &KafkaConsumer{
		sCtx: sCtx,
	}
	return kc
}
func (kc *KafkaConsumer) setupTracing(traceId string) context.Context {
	ctx := context.Background()
	if traceId != "" {
		if tid, err := trace.TraceIDFromHex(traceId); err == nil {
			spanCtx := trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid})
			ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
		}
	}
	return ctx
}
func (kc *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (kc *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (kc *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var data JobMsg
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logx.Errorf("消息反序列化失败, err: %v", err)
			session.MarkMessage(msg, "")
			continue
		}
		kc.handleMessage(msg, session, data)
	}
	return nil
}
func (kc *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession, data JobMsg) {
	ctx := kc.setupTracing(data.TraceID)
	l := logx.WithContext(ctx)
	l.Infof("开始处理任务 %d", data.AccountID)
	lockKey := fmt.Sprintf("lock:points:user:%d", data.UserID)
	isLocked, _ := kc.sCtx.RDB.SetNX(ctx, lockKey, "1", 10*time.Second).Result()
	if !isLocked {
		l.Infof("任务 %d 获取锁失败，稍后重试", data.AccountID)
		kc.enqueueDelay(data, "获取分布式锁失败")
		return
	}
	defer kc.sCtx.RDB.Del(ctx, lockKey)
	retryCountKey := fmt.Sprintf("retry:points:account:%d", data.AccountID)
	retryCount := kc.sCtx.RDB.Incr(ctx, retryCountKey)
	if retryCount.Err() != nil {
		l.Errorf("任务 %d 增加重试计数失败，err: %v", data.AccountID, retryCount.Err())
		kc.enqueueDelay(data, "增加重试计数失败")
		return
	}
	if retryCount.Val() > 3 {
		l.Errorf("任务 %d 超过最大重试次数，放弃处理", data.AccountID)
		kc.sCtx.RDB.Del(ctx, retryCountKey)
		kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, -1, "超过最大重试次数")
		return
	}
	kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, 1, "")
	user, err := kc.sCtx.PointsModel.FindUserById(ctx, data.UserID)
	if err != nil {
		l.Errorf("任务 %d 查询用户失败，err: %v", data.AccountID, err)
		kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, 0, "")
		kc.enqueueDelay(data, "查询用户失败")
		return
	}

	user.Score.Add(data.Amount)
	if user.Score.IsNegative() {
		l.Errorf("任务 %d 用户积分不足，无法扣除，err: %v", data.AccountID, err)
		kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, -1, "用户积分不足，无法扣除")
		return
	}
	err = kc.sCtx.PointsModel.UpdateScore(ctx, data.UserID, data.Amount)
	if err != nil {
		l.Errorf("任务 %d 更新用户积分失败，err: %v", data.AccountID, err)
		kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, 0, "")
		kc.enqueueDelay(data, "更新用户积分失败")
		return
	}
	err = kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, 2, "")
	if err != nil {
		l.Errorf("任务 %d 更新交易状态失败，err: %v", data.AccountID, err)
		return
	}
	kc.sCtx.PointsModel.UpdateTransactionStatus(ctx, strconv.FormatInt(data.AccountID, 10), data.UserID, 2, "")
	kc.sCtx.RDB.Del(ctx, retryCountKey)
	l.Infof("任务 %d 处理完成", data.AccountID)
}
func (kc *KafkaConsumer) enqueueDelay(data JobMsg, reason string) {
	delay := 5 * time.Second
	logx.Infof("任务 %d 进入时间轮，原因: %s, 将在 %v 后重试", data.AccountID, reason, delay)
	kc.sCtx.TimingWheel.AfterFunc(delay, func() {
		ctx := kc.setupTracing(data.TraceID)
		l := logx.WithContext(ctx)
		l.Infof("任务 %d 时间到达，重新入队", data.AccountID)
		val, _ := json.Marshal(data)
		mockMsg := &sarama.ConsumerMessage{
			Value: val,
		}
		kc.handleMessage(mockMsg, nil, data)
	})
}
