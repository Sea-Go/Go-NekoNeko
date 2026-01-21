// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"encoding/json"
	"fmt"

	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"
	"sea-try-go/service/article/rpc/articleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleLogic) CreateArticle(req *types.CreateArticleReq) (resp *types.CreateArticleResp, err error) {
	var authorId string
	if uid := l.ctx.Value("userId"); uid != nil {
		if idNum, ok := uid.(json.Number); ok {
			authorId = idNum.String()
		} else {
			authorId = fmt.Sprintf("%v", uid)
		}
	} else {
		authorId = "dev_test_user"
	}

	rpcResp, err := l.svcCtx.ArticleRpc.CreateArticle(l.ctx, &articleservice.CreateArticleRequest{
		Title:           req.Title,
		Brief:           &req.Brief,
		MarkdownContent: req.Content,
		CoverImageUrl:   &req.CoverImageUrl,
		ManualTypeTag:   req.ManualTypeTag,
		SecondaryTags:   req.SecondaryTags,
		AuthorId:        authorId,
	})

	if err != nil {
		return nil, err
	}

	return &types.CreateArticleResp{
		ArticleId: rpcResp.ArticleId,
	}, nil
}
