package mqs

import (
	"context"
	"encoding/json"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/svc"
	"time"
)

type UserPointsMsg struct {
	Uid        int64 `json:"uid"`
	AccountId  int64 `json:"accountId"`
	UserId     int64 `json:"userId"`
	Amount     int32 `json:"amount"`
	RetryTimes int32 `json:"retryTimes"`
}

type PointsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewPointsHandler(svcCtx *svc.ServiceContext) *PointsHandler {
	return &PointsHandler{svcCtx: svcCtx}
}

func (p *PointsHandler) Consume(ctx context.Context, key, value string) error {
	msg := &UserPointsMsg{}
	if err := json.Unmarshal([]byte(value), &msg); err != nil {
		return err
	}
	points, err := p.svcCtx.PointsModel.FindByAccountIdAndUserId(ctx, msg.AccountId, msg.UserId)
	if err != nil || points == nil {
		return err
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return nil
	}
	_, err = p.svcCtx.PointsModel.UpdateUserPoints(ctx, msg.UserId, msg.Amount)
	if err != nil {
		_, err := p.svcCtx.RetryDqPusherClient.Delay([]byte(value), time.Second*3)
		if err != nil {
			return err
		}
	}
	return nil
}
