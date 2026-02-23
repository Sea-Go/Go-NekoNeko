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

type UnfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnfollowLogic) Unfollow(in *pb.RelationReq) (*pb.BaseResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "Unfollow", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "Unfollow").Add(time.Since(start).Seconds())
	}()

	err := l.svcCtx.FollowModel.UnfollowUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		resultLabel = "sys_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "unfollow", "fail").Inc()
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbWrite, fmt.Errorf("UnfollowUser db err: %w", err))
		return &pb.BaseResp{Code: errmsg.ErrorDbWrite, Msg: errmsg.GetErrMsg(errmsg.ErrorDbWrite)}, err
	}
	metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "unfollow", "ok").Inc()
	return &pb.BaseResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success)}, nil
}
