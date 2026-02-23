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

type GetFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFollowerListLogic) GetFollowerList(in *pb.ListReq) (*pb.UserListResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "GetFollowerList", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "GetFollowerList").Add(time.Since(start).Seconds())
	}()

	ids, err := l.svcCtx.FollowModel.GetFollowerList(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		resultLabel = "sys_fail"
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbRead, fmt.Errorf("GetFollowerList db err: %w", err))
		return &pb.UserListResp{Code: errmsg.ErrorDbRead, Msg: errmsg.GetErrMsg(errmsg.ErrorDbRead)}, err
	}

	metrics.FollowListSizeGaugeMetric.WithLabelValues("follow_list", "follower").Set(float64(len(ids)))
	return &pb.UserListResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success), UserIds: ids}, nil
}
