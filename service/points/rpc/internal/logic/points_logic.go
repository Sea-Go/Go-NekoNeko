package logic

import (
	"context"
	"encoding/json"
	"errors"
	"sea-try-go/service/common/dberr"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/mq"
	"sea-try-go/service/points/rpc/internal/svc"
	"sea-try-go/service/points/rpc/pb"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
)

type PointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointsLogic {
	return &PointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PointsLogic) Points(in *__.PointsReq) (*__.PointsResp, error) {
	amount, _ := decimal.NewFromString(in.Amount)
	teaceId := ""
	spanContext := trace.SpanContextFromContext(l.ctx)
	if spanContext.HasTraceID() {
		teaceId = spanContext.TraceID().String()
	}
	txRecord := &model.Transaction{
		AccountId: strconv.FormatInt(in.AccountId, 10),
		UserId:    in.UserId,
		Amount:    amount,
		Tracing:   teaceId,
	}
	err := l.svcCtx.PointsModel.CreateTransaction(l.ctx, txRecord)
	if err != nil {
		if dberr.IsDuplicateKeyError(err) {
			return &__.PointsResp{Success: true, Message: "请求已存在，正在处理中"}, nil
		}
		logx.Errorf("创建交易记录失败, err: %v", err)
		return nil, err
	}
	kafkaMsg := mq.JobMsg{
		AccountID: in.AccountId,
		UserID:    in.UserId,
		Amount:    amount,
		TraceID:   teaceId,
	}
	msgBytes, _ := json.Marshal(kafkaMsg)
	kMsg := &sarama.ProducerMessage{
		Topic: l.svcCtx.Config.KafkaConf.Topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(in.UserId, 10)),
		Value: sarama.ByteEncoder(msgBytes),
	}
	if _, _, err := l.svcCtx.KafKaProducer.SendMessage(kMsg); err != nil {
		logx.Errorf("发送Kafka消息失败, err: %v", err)
		err := l.svcCtx.PointsModel.UpdateTransactionStatus(l.ctx, txRecord.AccountId, txRecord.UserId, -1, "发送Kafka消息失败")
		if err != nil {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
		}
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}

	return &__.PointsResp{
		Success: true,
		Message: "请求已接受，正在处理中",
	}, nil
}
