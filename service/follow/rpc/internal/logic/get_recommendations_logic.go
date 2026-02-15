package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecommendationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendationsLogic {
	return &GetRecommendationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetRecommendationsLogic) GetRecommendations(in *pb.ListReq) (*pb.RecommendResp, error) {
	recs, err := l.svcCtx.FollowModel.GetRecommendations(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		l.Logger.Errorf("GetRecommendations db err: %v", err)
		return &pb.RecommendResp{Code: 500, Msg: "DB Error"}, err
	}

	// 将 Model 返回的内部结构体 转换为 pb (契约) 定义的返回格式
	var pbUsers []*pb.RecommendResp_RecommendUser
	for _, rec := range recs {
		pbUsers = append(pbUsers, &pb.RecommendResp_RecommendUser{
			TargetId:    rec.TargetId,
			MutualScore: int32(rec.MutualScore),
		})
	}

	return &pb.RecommendResp{Code: 0, Msg: "success", Users: pbUsers}, nil
}
