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

type GetBlockListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockListLogic {
	return &GetBlockListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetBlockListLogic) GetBlockList(in *pb.ListReq) (*pb.UserListResp, error) {
	start := time.Now()
	resultLabel := "ok"
	defer func() {
		metrics.FollowRequestCounterMetric.WithLabelValues("follow_rpc", "GetBlockList", resultLabel).Inc()
		metrics.FollowRequestSecondsCounterMetric.WithLabelValues("follow_rpc", "GetBlockList").Add(time.Since(start).Seconds())
	}()

	ids, err := l.svcCtx.FollowModel.GetBlockList(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		resultLabel = "sys_fail"
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbRead, fmt.Errorf("GetBlockList db err: %w", err))
		return &pb.UserListResp{Code: errmsg.ErrorDbRead, Msg: errmsg.GetErrMsg(errmsg.ErrorDbRead)}, err
	}

	metrics.FollowListSizeGaugeMetric.WithLabelValues("follow_list", "blocked").Set(float64(len(ids)))
	return &pb.UserListResp{Code: errmsg.Success, Msg: errmsg.GetErrMsg(errmsg.Success), UserIds: ids}, nil
}
