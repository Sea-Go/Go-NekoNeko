package mqs

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/svc"
)

type DelayHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDelayHandler(svcCtx *svc.ServiceContext) *DelayHandler {
	return &DelayHandler{svcCtx: svcCtx}
}

func (h *DelayHandler) Consume(body []byte) {
	ctx := context.Background()
	var uid int64
	if err := json.Unmarshal(body, &uid); err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorJsonUnmarshal, err)
		return
	}

	points, err := h.svcCtx.PointsModel.FindOneByUid(ctx, uid)
	if err != nil || points == nil {
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, err)
		}
		return
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return
	}
	// 15分钟延时到达仍未完成，标记为超时失败
	logger.LogBusinessErr(ctx, errmsg.ErrorPointsTimeout,
		fmt.Errorf("uid=%d 积分处理超时，已超过15分钟", uid),
		logger.WithUserID(fmt.Sprintf("%d", points.UserId)),
	)
	h.svcCtx.PointsModel.UpdateStatusByUid(ctx, uid, model.StatusFailed)
}
