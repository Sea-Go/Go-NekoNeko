package mqs

import (
	"context"
	"encoding/json"
	kqtypes "sea-try-go/service/comment/common/types"
	"sea-try-go/service/comment/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *AuditConsumer {
	return &AuditConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuditConsumer) Consume(ctx context.Context, key, val string) error {
	var msg kqtypes.CommentKafkaMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("Kafka消息解析失败:%v", err)
		return nil
	}
	isSensitive, hitWord := l.svcCtx.SensitiveFilter.Match(msg.Content)
	status := 1
	if isSensitive {
		l.Infof("发现违规词 [%s],评论ID :%d 被打入审核池", hitWord, msg.CommentId)
		status = 2
	}
	err := l.svcCtx.CommentModel.InsertCommentTx(ctx, msg, status)
	if err != nil {
		l.Errorf("[写库失败] 评论 ID: %d 写入 PG 失败: %v", msg.CommentId, err)
		return err
	}
	l.Infof("评论流转完成, ID: %d, 最终状态: %d", msg.CommentId, status)
	return nil
}
