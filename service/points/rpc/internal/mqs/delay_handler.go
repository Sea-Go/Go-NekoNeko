package mqs

import (
	"context"
	"encoding/json"
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
	var uid int64
	if err := json.Unmarshal(body, &uid); err != nil {
		return
	}

	ctx := context.Background()
	points, err := h.svcCtx.PointsModel.FindOneByUid(ctx, uid)
	if err != nil || points == nil {
		return
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return
	}
	// trace 链路追踪
}
