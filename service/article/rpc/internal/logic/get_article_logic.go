package logic

import (
	"context"
	"io"

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

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
		l.Logger.Errorf("GetArticle db error: %v", err)
		return nil, err
	}

	if in.IncrView {
		if err := l.svcCtx.ArticleRepo.IncrViewCount(l.ctx, in.ArticleId); err != nil {
			l.Logger.Errorf("IncrViewCount error: %v", err)
		}
		article.ViewCount++
	}

	object, err := l.svcCtx.MinioClient.GetObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, article.Content, minio.GetObjectOptions{})
	if err != nil {
		l.Logger.Errorf("MinIO GetObject error: %v", err)
		return nil, err
	}
	defer object.Close()

	contentBytes, err := io.ReadAll(object)
	if err != nil {
		l.Logger.Errorf("Read MinIO object error: %v", err)
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
