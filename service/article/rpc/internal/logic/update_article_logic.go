package logic

import (
	"context"
	"strings"

	"fmt"

	"sea-try-go/service/article/common/errmsg"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"
	"sea-try-go/service/common/logger"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateArticleLogic) UpdateArticle(in *__.UpdateArticleRequest) (*__.UpdateArticleResponse, error) {
	article, err := l.svcCtx.ArticleRepo.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err, logger.WithArticleID(in.ArticleId))
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	if in.Title != nil {
		article.Title = *in.Title
	}
	if in.Brief != nil {
		article.Brief = *in.Brief
	}
	if in.MarkdownContent != nil {
		objectName := article.Content
		if objectName == "" {
			objectName = fmt.Sprintf("%s%s.md", l.svcCtx.Config.MinIO.ArticlePath, article.ID)
			article.Content = objectName
		}

		contentType := "text/markdown"
		reader := strings.NewReader(*in.MarkdownContent)
		_, err = l.svcCtx.MinioClient.PutObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, objectName,
			reader, int64(len(*in.MarkdownContent)), minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			logger.LogBusinessErr(l.ctx, errmsg.Error, fmt.Errorf("update minio content failed: %w", err), logger.WithArticleID(in.ArticleId))
			return nil, err
		}
	}
	if in.CoverImageUrl != nil {
		article.CoverImageURL = *in.CoverImageUrl
	}
	if in.ManualTypeTag != nil {
		article.ManualTypeTag = *in.ManualTypeTag
	}
	if len(in.SecondaryTags) > 0 {
		article.SecondaryTags = in.SecondaryTags
	}
	if in.Status != nil && *in.Status != __.ArticleStatus_ARTICLE_STATUS_UNSPECIFIED {
		article.Status = int32(in.Status.Number())
	}

	if err := l.svcCtx.ArticleRepo.Update(l.ctx, article); err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err, logger.WithArticleID(in.ArticleId))
		return nil, err
	}

	return &__.UpdateArticleResponse{
		Success: true,
	}, nil
}
