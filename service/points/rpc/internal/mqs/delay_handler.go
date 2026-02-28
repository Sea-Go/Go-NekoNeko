package mqs

import (
	"context"
	"encoding/json"
	"fmt"

	"sea-try-go/service/common/logger"
	"sea-try-go/service/points/rpc/internal/metrics"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/svc"
	"sea-try-go/service/user/common/errmsg"
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
		metrics.PointsKafkaErrorCounterMetric.WithLabelValues("points_mq", "delay_consume", "json_unmarshal").Inc()
		logger.LogBusinessErr(ctx, errmsg.ErrorJsonUnmarshal, err)
		return
	}

	points, err := h.svcCtx.PointsModel.FindOneByUid(ctx, uid)
	if err != nil || points == nil {
		if err != nil {
			metrics.PointsKafkaErrorCounterMetric.WithLabelValues("points_mq", "delay_consume", "db_select").Inc()
			logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, err)
		}
		return
	}
	if points.Status == model.StatusSuccess || points.Status == model.StatusFailed {
		return
	}

	metrics.PointsKafkaErrorCounterMetric.WithLabelValues("points_mq", "delay_consume", "timeout").Inc()
	logger.LogBusinessErr(ctx, errmsg.ErrorPointsTimeout,
		fmt.Errorf("uid=%d points process timeout after 15 minutes", uid),
		logger.WithUserID(fmt.Sprintf("%d", points.UserId)),
	)
	if err = h.svcCtx.PointsModel.UpdateStatusByUid(ctx, uid, model.StatusFailed); err != nil {
		metrics.PointsKafkaErrorCounterMetric.WithLabelValues("points_mq", "delay_consume", "status_update").Inc()
		logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, err, logger.WithUserID(fmt.Sprintf("%d", points.UserId)))
	}
}
