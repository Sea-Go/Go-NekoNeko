package mqs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sea-try-go/service/article/common/errmsg"
	"sea-try-go/service/article/rpc/internal/svc"
	pb "sea-try-go/service/article/rpc/pb"
	"sea-try-go/service/common/logger"

	green "github.com/alibabacloud-go/green-20220302/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleConsumer {
	return &ArticleConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleConsumer) Consume(ctx context.Context, key, val string) error {
	logger.LogInfo(ctx, fmt.Sprintf("DataClean Service Consuming: %s", val))

	var msg struct {
		ArticleId string `json:"article_id"`
		AuthorId  string `json:"author_id"`
		Content   string `json:"content"`
	}

	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("unmarshal error: %w", err))
		return nil
	}

	article, err := l.svcCtx.ArticleRepo.FindOne(ctx, msg.ArticleId)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, fmt.Errorf("failed to find article %s: %w", msg.ArticleId, err), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
		return err
	}

	// Idempotency check: skip if not in pending status (3)
	if article.Status != int32(pb.ArticleStatus_REVIEWING) {
		logger.LogInfo(ctx, fmt.Sprintf("Article %s status is %d, skipping duplicate processing.", msg.ArticleId, article.Status))
		return nil
	}

	if l.svcCtx.GreenClient == nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("GreenClient is not initialized"), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
		return fmt.Errorf("GreenClient not initialized")
	}

	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": msg.Content,
		},
	)
	request := green.TextModerationRequest{
		Service:           tea.String("comment_detection"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, err := l.svcCtx.GreenClient.TextModeration(&request)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("AliGreen API error: %w", err), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
		return err
	}

	statusCode := tea.IntValue(tea.ToInt(result.StatusCode))
	if statusCode == http.StatusOK {
		textModerationResponse := result.Body
		if tea.IntValue(tea.ToInt(textModerationResponse.Code)) == 200 {
			textModerationResponseData := textModerationResponse.Data
			reason := tea.StringValue(textModerationResponseData.Reason)
			labels := tea.StringValue(textModerationResponseData.Labels)

			if len(reason) > 0 || len(labels) > 0 {
				logger.LogInfo(ctx, fmt.Sprintf("Article %s RISK DETECTED! Reason: %s, Labels: %s", msg.ArticleId, reason, labels), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
				article.Status = int32(pb.ArticleStatus_REJECTED)
				if err := l.svcCtx.ArticleRepo.Update(ctx, article); err != nil {
					logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, fmt.Errorf("failed to update article status to Rejected: %w", err), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
					return err
				}
			} else {
				logger.LogInfo(ctx, fmt.Sprintf("Article %s passed safety check.", msg.ArticleId), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))

				bucketName := l.svcCtx.Config.MinIO.BucketName
				objectName := fmt.Sprintf("%s.md", msg.ArticleId)
				reader := strings.NewReader(msg.Content)

				_, err = l.svcCtx.MinioClient.PutObject(ctx, bucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{
					ContentType: "text/markdown",
				})
				if err != nil {
					logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("failed to upload to MinIO: %w", err), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
					return err
				}
				
				logger.LogInfo(ctx, fmt.Sprintf("Article %s uploaded to MinIO bucket %s", msg.ArticleId, bucketName), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))

				article.Status = int32(pb.ArticleStatus_PUBLISHED)
				if err := l.svcCtx.ArticleRepo.Update(ctx, article); err != nil {
					logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, fmt.Errorf("failed to update article status to Published: %w", err), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
					return err
				}
			}
		} else {
			logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("AliGreen response code error: %v", textModerationResponse.Code), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
		}
	} else {
		logger.LogBusinessErr(ctx, errmsg.ErrorServerCommon, fmt.Errorf("AliGreen http status error: %v", statusCode), logger.WithArticleID(msg.ArticleId), logger.WithUserID(msg.AuthorId))
	}

	return nil
}
