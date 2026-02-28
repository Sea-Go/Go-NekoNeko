package logic

import (
```
	"context"
	"fmt" // Added for fmt.Errorf

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/minio/minio-go/v7" // 引入 minio 包
)

```

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
```
	article, err := l.svcCtx.ArticleRepo.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		l.Logger.Errorf("DeleteArticle FindOne error: %v", err)
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	// 1. 从 MinIO 删除对应的 Markdown 内容
	if article.Content != "" { // 确保有内容路径才删除
		err = l.svcCtx.MinioClient.RemoveObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, article.Content, minio.RemoveObjectOptions{})
		if err != nil {
			l.Logger.Errorf("DeleteArticle RemoveObject from MinIO error: %v, object: %s", err, article.Content)
			// 这里可以选择是否中断删除操作。通常，即使 MinIO 删除失败，数据库记录也应该被删除。
			// 但为了数据一致性，这里选择返回错误，让用户重试或手动处理。
			return nil, err
		}
	}

	// 2. 删除数据库记录
	if err := l.svcCtx.ArticleRepo.Delete(l.ctx, in.ArticleId); err != nil {
		l.Logger.Errorf("DeleteArticle Delete error: %v", err)
		return nil, err
	}

```
	return &__.DeleteArticleResponse{
		Success: true,
	}, nil
}
