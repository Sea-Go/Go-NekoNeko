package logic

import (
	"context"
	"fmt"
	"time"

	"sea-try-go/service/common/logger"
	"sea-try-go/service/follow/common/errmsg"
	"sea-try-go/service/follow/rpc/internal/metrics"
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
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "GetRecommendations", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "GetRecommendations").Add(time.Since(start).Seconds())
	}()

	recs, err := l.svcCtx.FollowModel.GetRecommendations(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		resultLabel = "sys_fail"
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbRead, fmt.Errorf("GetRecommendations db err: %w", err))
		return &pb.RecommendResp{Code: errmsg.ErrorDbRead, Msg: errmsg.GetErrMsg(errmsg.ErrorDbRead)}, err
	}

	// 将 Model 返回的内部结构体 转换为 pb (契约) 定义的返回格式
	var pbUsers []*pb.RecommendResp_RecommendUser
	for _, rec := range recs {
		pbUsers = append(pbUsers, &pb.RecommendResp_RecommendUser{
			TargetId:    rec.TargetId,
			MutualScore: int32(rec.MutualScore),
		})
	}

	metrics.FollowListSizeGaugeMetric.WithLabelValues("follow_list", "recommendation").Set(float64(len(pbUsers)))
	return &pb.RecommendResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success), Users: pbUsers}, nil
}
