package logic

import (
	"context"
	"fmt"

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleLogic) DeleteArticle(in *__.DeleteArticleRequest) (*__.DeleteArticleResponse, error) {
	article, err := l.svcCtx.ArticleRepo.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		l.Logger.Errorf("DeleteArticle FindOne error: %v", err)
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	if article.Content != "" {
		err = l.svcCtx.MinioClient.RemoveObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, article.Content, minio.RemoveObjectOptions{})
		if err != nil {
			l.Logger.Errorf("DeleteArticle RemoveObject from MinIO error: %v, object: %s", err, article.Content)
			return nil, err
		}
	}

	if err := l.svcCtx.ArticleRepo.Delete(l.ctx, in.ArticleId); err != nil {
		l.Logger.Errorf("DeleteArticle Delete error: %v", err)
		return nil, err
	}

	return &__.DeleteArticleResponse{
		Success: true,
	}, nil
}
