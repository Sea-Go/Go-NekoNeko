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

type UnblockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnblockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockLogic {
	return &UnblockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnblockLogic) Unblock(in *pb.RelationReq) (*pb.BaseResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "Unblock", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "Unblock").Add(time.Since(start).Seconds())
	}()

	err := l.svcCtx.FollowModel.UnblockUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		resultLabel = "sys_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "unblock", "fail").Inc()
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbWrite, fmt.Errorf("UnblockUser db err: %w", err))
		return &pb.BaseResp{Code: errmsg.ErrorDbWrite, Msg: errmsg.GetErrMsg(errmsg.ErrorDbWrite)}, err
	}
	metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "unblock", "ok").Inc()
	return &pb.BaseResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success)}, nil
}
