package mqs

import (
	"context"
	"encoding/json"
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
		//return err
		return
	}
	ctx := context.Background()
	points, err := h.svcCtx.PointsModel.FindByAccountIdAndUserId(ctx, msg.AccountId, msg.UserId)
	if err != nil || points == nil {
		//return err
		return
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		//return nil
	}
	if msg.RetryTimes > 3 {
		err = h.svcCtx.PointsModel.UpdateStatusByUid(ctx, msg.Uid, model.StatusFailed)
		if err != nil {
			//return err
			return
		}
	}
	_, err = h.svcCtx.PointsModel.UpdateUserPoints(ctx, msg.UserId, msg.Amount)
	if err != nil {
		msg.RetryTimes += 1
		bytes, err := json.Marshal(msg)
		if err != nil {
			//return err
			return
		}
		h.svcCtx.RetryDqPusherClient.Delay(bytes, time.Second*3)

	}
}
