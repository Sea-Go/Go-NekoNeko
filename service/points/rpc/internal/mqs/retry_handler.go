package mqs

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/svc"
	"time"
)

type RetryHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRetryHandler(svcCtx *svc.ServiceContext) *RetryHandler {
	return &RetryHandler{svcCtx: svcCtx}
}

func (h *RetryHandler) Consume(body []byte) {
	msg := &UserPointsMsg{}
	if err := json.Unmarshal(body, msg); err != nil {
		logger.LogBusinessErr(context.Background(), errmsg.ErrorJsonUnmarshal, err)
		return
	}
	ctx := context.Background()
	points, err := h.svcCtx.PointsModel.FindByAccountIdAndUserId(ctx, msg.AccountId, msg.UserId)
	if err != nil || points == nil {
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, err)
		}
		return
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return
	}
	if msg.RetryTimes > 3 {
		logger.LogBusinessErr(ctx, errmsg.ErrorPointsRetryExceeded, fmt.Errorf("uid=%d 重试次数超限: %d", msg.Uid, msg.RetryTimes), logger.WithUserID(fmt.Sprintf("%d", msg.UserId)))
		err = h.svcCtx.PointsModel.UpdateStatusByUid(ctx, msg.Uid, model.StatusFailed)
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, err)
		}
		return
	}
	ok, err := h.svcCtx.PointsModel.UpdateUserPoints(ctx, msg.UserId, msg.Amount)
	if err != nil {
		msg.RetryTimes += 1
		bytes, err := json.Marshal(msg)
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorJsonMarshal, err)
			return
		}
		h.svcCtx.RetryDqPusherClient.Delay(bytes, time.Second*3)
		return
	}
	if !ok {
		// 积分余量不足，标记为失败
		logger.LogBusinessErr(ctx, errmsg.ErrorPointsInsufficient, fmt.Errorf("userId=%d 积分余量不足", msg.UserId), logger.WithUserID(fmt.Sprintf("%d", msg.UserId)))
		h.svcCtx.PointsModel.UpdateStatusByUid(ctx, msg.Uid, model.StatusFailed)
	}
}
