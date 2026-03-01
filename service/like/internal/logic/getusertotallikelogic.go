package logic

import (
	"context"
	"fmt"
	"strconv"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTotalLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserTotalLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTotalLikeLogic {
	return &GetUserTotalLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserTotalLikeLogic) GetUserTotalLike(in *pb.GetUserTotalLikeReq) (*pb.GetUserTotalLikeResp, error) {
	redisKey := fmt.Sprintf("user_total_like:%d", in.UserId)
	val, err := l.svcCtx.BizRedis.GetCtx(l.ctx, redisKey)
	if err == nil && val != "" {
		count, _ := strconv.ParseInt(val, 10, 64)
		return &pb.GetUserTotalLikeResp{
			TotalLikeCount: count,
		}, nil
	}
	totalCount, err := l.svcCtx.LikeModel.GetTotalLikeCount(in.UserId)
	if err != nil {
		return nil, err
	}

	ttl := l.svcCtx.Config.Storage.Redis.CacheTTL
	err = l.svcCtx.BizRedis.SetexCtx(l.ctx, redisKey, strconv.FormatInt(totalCount, 10), ttl)

	//Redis可能出问题了
	if err != nil {
		return nil, err
	}
	return &pb.GetUserTotalLikeResp{
		TotalLikeCount: totalCount,
	}, nil
}
