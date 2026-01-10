package logic

import (
	"context"

	"sea-try-go/service/favorite/rpc/internal/svc"
	"sea-try-go/service/favorite/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(in *__.FavoriteListReq) (*__.FavoriteListResp, error) {
	res, err := l.svcCtx.FavoriteRepo.GetArtocleIdListByUserId(in.UserId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logx.Errorf("GetFavoriteListLogic.GetArtocleIdListByUserId(),userId:%d,error: %v", in.UserId, err)
		return nil, err
	} else if err == gorm.ErrRecordNotFound || res == nil {
		return &__.FavoriteListResp{}, nil
	}
	var finRes __.FavoriteListResp
	finRes.List = make([]*__.FavoriteItem, 0)
	for _, i := range *res {
		finRes.List = append(finRes.List, &__.FavoriteItem{
			ArticleId: i.ArticleID,
		})
		finRes.Total++
	}
	return &finRes, nil
}
