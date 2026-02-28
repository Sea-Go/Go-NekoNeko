package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"sea-try-go/service/article/common/errmsg"
	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/common/snowflake"
	"strings"

	"github.com/minio/minio-go/v7"
)

type CreateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateArticleLogic) CreateArticle(in *__.CreateArticleRequest) (*__.CreateArticleResponse, error) {
	idInt, err := snowflake.GetID()
	if err != nil {
		return nil, err
	}
	articleId := fmt.Sprintf("%d", idInt)

	objectName := fmt.Sprintf("%s%s.md", l.svcCtx.Config.MinIO.ArticlePath, articleId)
	contentType := "text/markdown"
	reader := strings.NewReader(in.MarkdownContent)
	_, err = l.svcCtx.MinioClient.PutObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, objectName,
		reader, int64(len(in.MarkdownContent)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.Error, err, logger.WithArticleID(articleId))
		return nil, err
	}

	newArticle := &model.Article{
		ID:            articleId,
		Title:         in.Title,
		Brief:         *in.Brief,
		Content:       objectName, // 这里存的是 MinIO 的路径，而不是原文
		CoverImageURL: *in.CoverImageUrl,
		ManualTypeTag: in.ManualTypeTag,
		SecondaryTags: model.StringArray(in.SecondaryTags),
		AuthorID:      in.AuthorId,
		Status:        int32(__.ArticleStatus_REVIEWING),
	}

	if err := l.svcCtx.ArticleRepo.Insert(l.ctx, newArticle); err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err, logger.WithArticleID(articleId), logger.WithUserID(in.AuthorId))
		return nil, err
	}

	msg := struct {
		ArticleId   string `json:"article_id"`
		AuthorId    string `json:"author_id"`
		ContentPath string `json:"content_path"`
	}{
		ArticleId:   articleId,
		AuthorId:    in.AuthorId,
		ContentPath: objectName,
	}

	msgBytes, _ := json.Marshal(msg)
	if err := l.svcCtx.KqPusher.Push(l.ctx, string(msgBytes)); err != nil {
		err = fmt.Errorf("kafka push failed, payload: %s, error: %w", string(msgBytes), err)
		logger.LogBusinessErr(l.ctx, errmsg.Error, err, logger.WithArticleID(articleId), logger.WithUserID(in.AuthorId))
	}

	return &__.CreateArticleResponse{
		ArticleId: articleId,
	}, nil
}
