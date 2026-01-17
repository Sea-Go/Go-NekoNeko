package logic

import (
	"context"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeArticleLogic) LikeArticle(in *__.LikeArticleReq) (*__.LikeArticleResp, error) {
	// todo: add your logic here and delete this line

	return &__.LikeArticleResp{
		LikeCount: 0, // 暂时填0，你需要做类型转换
	}, nil
}
