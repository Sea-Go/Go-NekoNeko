package logic

import (
	"context"
	"fmt"
	"io"

	"sea-try-go/service/article/common/errmsg"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"
	"sea-try-go/service/common/logger"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleLogic) GetArticle(in *__.GetArticleRequest) (*__.GetArticleResponse, error) {
	article, err := l.svcCtx.ArticleRepo.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err, logger.WithArticleID(in.ArticleId))
		return nil, err
	}

	if in.IncrView {
		if err := l.svcCtx.ArticleRepo.IncrViewCount(l.ctx, in.ArticleId); err != nil {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err, logger.WithArticleID(in.ArticleId))
		}
		article.ViewCount++
	}

	object, err := l.svcCtx.MinioClient.GetObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, article.Content, minio.GetObjectOptions{})
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorMinioDownload, fmt.Errorf("minio get object failed: %w", err), logger.WithArticleID(in.ArticleId))
		return nil, err
	}
	defer object.Close()

	contentBytes, err := io.ReadAll(object)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorMinioDownload, fmt.Errorf("read minio object failed: %w", err), logger.WithArticleID(in.ArticleId))
		return nil, err
	}

	return &__.GetArticleResponse{
		Article: &__.Article{
			Id:              article.ID,
			Title:           article.Title,
			Brief:           article.Brief,
			MarkdownContent: string(contentBytes),
			CoverImageUrl:   article.CoverImageURL,
			ManualTypeTag:   article.ManualTypeTag,
			SecondaryTags:   article.SecondaryTags,
			AuthorId:        article.AuthorID,
			CreateTime:      article.CreatedAt.UnixMilli(),
			UpdateTime:      article.UpdatedAt.UnixMilli(),
			Status:          __.ArticleStatus(article.Status),
			ViewCount:       article.ViewCount,
			LikeCount:       article.LikeCount,
			CommentCount:    article.CommentCount,
		},
	}, nil
}
