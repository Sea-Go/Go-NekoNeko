package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/common/utils"
	"sea-try-go/service/article/rpc/internal/model"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"
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

	newArticle := &model.Article{
		ID:            articleId,
		Title:         in.Title,
		Brief:         *in.Brief,
		Content:       in.MarkdownContent,
		CoverImageURL: *in.CoverImageUrl,
		ManualTypeTag: in.ManualTypeTag,
		SecondaryTags: model.StringArray(in.SecondaryTags),
		AuthorID:      in.AuthorId,
		Status:        3,
	}

	if err := l.svcCtx.ArticleRepo.Db.WithContext(l.ctx).Create(newArticle).Error; err != nil {
		l.Logger.Errorf("CreateArticle db error: %v", err)
		return nil, err
	}

	msg := struct {
		ArticleId string `json:"article_id"`
		AuthorId  string `json:"author_id"`
		Action    string `json:"action"`
		Timestamp int64  `json:"timestamp"`
	}{
		ArticleId: articleId,
		AuthorId:  in.AuthorId,
		Action:    "create",
		Timestamp: time.Now().Unix(),
	}

	msgBytes, _ := json.Marshal(msg)
	if err := l.svcCtx.KqPusher.Push(l.ctx, string(msgBytes)); err != nil {
		l.Logger.Errorf("Failed to push article creation event to Kafka: %v, payload: %s", err, string(msgBytes))
	}

	return &__.CreateArticleResponse{
		ArticleId: articleId,
	}, nil
}
