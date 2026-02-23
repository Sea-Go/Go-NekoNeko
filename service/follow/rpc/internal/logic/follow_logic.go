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

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FollowLogic) Follow(in *pb.RelationReq) (*pb.BaseResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "Follow", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "Follow").Add(time.Since(start).Seconds())
	}()

	// 1. 业务特判
	if in.UserId == in.TargetId {
		resultLabel = "biz_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "follow", "fail").Inc()
		return &pb.BaseResp{Code: errmsg.ErrorCannotFollowSelf, Msg: errmsg.GetErrMsg(errmsg.ErrorCannotFollowSelf)}, nil
	}

	// 2. 调用 Model 层干活
	err := l.svcCtx.FollowModel.FollowUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		resultLabel = "sys_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "follow", "fail").Inc()
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbWrite, fmt.Errorf("FollowUser db err: %w", err))
		return &pb.BaseResp{Code: errmsg.ErrorDbWrite, Msg: errmsg.GetErrMsg(errmsg.ErrorDbWrite)}, err
	}

	metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "follow", "ok").Inc()
	return &pb.BaseResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success)}, nil
}
