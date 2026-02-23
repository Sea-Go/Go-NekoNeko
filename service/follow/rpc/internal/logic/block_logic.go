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

type BlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockLogic {
	return &BlockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BlockLogic) Block(in *pb.RelationReq) (*pb.BaseResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "Block", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "Block").Add(time.Since(start).Seconds())
	}()

	if in.UserId == in.TargetId {
		resultLabel = "biz_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "block", "fail").Inc()
		return &pb.BaseResp{Code: errmsg.ErrorCannotBlockSelf, Msg: errmsg.GetErrMsg(errmsg.ErrorCannotBlockSelf)}, nil
	}

	err := l.svcCtx.FollowModel.BlockUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		resultLabel = "sys_fail"
		metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "block", "fail").Inc()
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbWrite, fmt.Errorf("BlockUser db err: %w", err))
		return &pb.BaseResp{Code: errmsg.ErrorDbWrite, Msg: errmsg.GetErrMsg(errmsg.ErrorDbWrite)}, err
	}
	metrics.FollowRelationCounterMetric.WithLabelValues("follow_relation", "block", "ok").Inc()
	return &pb.BaseResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success)}, nil
}
