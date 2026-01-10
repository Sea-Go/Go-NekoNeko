package logic

import (
	"context"

	"sea-try-go/service/favorite/rpc/internal/svc"
	"sea-try-go/service/favorite/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteActionLogic) FavoriteAction(in *__.FavoriteActionReq) (*__.FavoriteActionResp, error) {
	userId := in.UserId
	articleId := in.ArticleId
	actionType := in.ActionType
	if actionType == 0 {
		err := l.svcCtx.FavoriteRepo.Delete(userId, articleId)
		if err != nil && err != gorm.ErrRecordNotFound {
			logx.Errorf(
				"favorite delete failed, user=%d article=%d err=%v",
				userId, articleId, err,
			)
			return nil, err
		}
	} else {
		err := l.svcCtx.FavoriteRepo.Insert(userId, articleId)
		if err != nil {
			logx.Errorf(
				"favorite insert failed, user=%d article=%d err=%v",
				userId, articleId, err,
			)
			return nil, err
		}
	}
	return &__.FavoriteActionResp{Success: true}, nil
}
