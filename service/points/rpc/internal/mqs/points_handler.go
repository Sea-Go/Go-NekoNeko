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
		logger.LogBusinessErr(ctx, errmsg.ErrorJsonUnmarshal, err)
		return err
	}
	points, err := p.svcCtx.PointsModel.FindByAccountIdAndUserId(ctx, msg.AccountId, msg.UserId)
	if err != nil || points == nil {
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, err)
		}
		return err
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return nil
	}
	ok, err := p.svcCtx.PointsModel.UpdateUserPoints(ctx, msg.UserId, msg.Amount)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, err, logger.WithUserID(fmt.Sprintf("%d", msg.UserId)))
		_, err := p.svcCtx.RetryDqPusherClient.Delay([]byte(value), time.Second*3)
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDelayMsg, err)
			return err
		}
		return nil
	}
	if !ok {
		// 积分余量不足，标记为失败
		err = p.svcCtx.PointsModel.UpdateStatusByUid(ctx, msg.Uid, model.StatusFailed)
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, err, logger.WithUserID(fmt.Sprintf("%d", msg.UserId)))
			return err
		}
	}
	return nil
}
