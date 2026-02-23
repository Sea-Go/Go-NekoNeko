package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/points/rpc/internal/svc"
	"strconv"
	"time"
)

// UserPointsMsg Kafka 消息结构
type UserPointsMsg struct {
	Uid        int64 `json:"uid"`
	AccountId  int64 `json:"accountId"`
	UserId     int64 `json:"userId"`
	Amount     int32 `json:"amount"`
	RetryTimes int32 `json:"retryTimes"`
}

// PointsResult 统一的积分操作返回结构
type PointsResult struct {
	Success bool
	Message string
}

// ProcessPoints 公共积分处理流程，add 和 dec 共用
// amount: 正数表示加积分，负数表示扣积分
func ProcessPoints(ctx context.Context, svcCtx *svc.ServiceContext, userId int64, accountId int64, amount int32) (*PointsResult, error) {
	// logger := logx.WithContext(ctx)

	// 使用雪花算法生成唯一 Uid
	uid, err := snowflake.GetID()
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorSnowflakeID, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
		return nil, err
	}

	// Step 1: 检查是否有正在处理的任务
	hasProcessing, err := svcCtx.PointsModel.HasProcessingByUserId(ctx, userId)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorDbSelect, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
		return nil, err
	}
	if hasProcessing {
		// 有人排队 -> 插入 Queued 记录并返回
		err := svcCtx.PointsModel.Insert(ctx, &model.Points{
			Uid:       uid,
			UserId:    userId,
			AccountId: accountId,
			Amount:    amount,
			Status:    model.StatusQueued,
		})
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbInsert, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
			return nil, err
		}
		return &PointsResult{Success: true, Message: "任务已排队"}, nil
	}

	// 无人排队 -> 开始处理流程
	// Step 2: 开启事务
	tx := svcCtx.PointsModel.BeginTransaction()
	pointsLog := &model.Points{
		Uid:       uid, // 雪花算法生成的唯一ID
		AccountId: accountId,
		UserId:    userId,
		Amount:    amount,
		Status:    model.StatusProcessing,
	}
	if err := tx.Create(pointsLog).Error; err != nil {
		tx.Rollback()
		logger.LogBusinessErr(ctx, errmsg.ErrorDbInsert, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
		return nil, err
	}

	// Step 4: 发送延时消息 (带重试)
	for i := 0; i < 3; i++ {
		if err = sendDelayCheck(ctx, svcCtx, uid); err == nil {
			break
		}
		if i == 2 {
			tx.Rollback()
			logger.LogBusinessErr(ctx, errmsg.ErrorDelayMsg, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
			return nil, errors.New("发送延时消息失败")
		}
	}

	// Step 5: 提交事务
	tx.Commit()

	// Step 6: Double Check
	hasOther, _ := svcCtx.PointsModel.HasOtherProcessingByUserId(ctx, userId, pointsLog.Uid)
	if hasOther {
		// 发现冲突，降级为 Queued
		logger.LogInfo(ctx, "并发冲突，任务降级为排队", logger.WithUserID(fmt.Sprintf("%d", userId)))
		err := svcCtx.PointsModel.UpdateStatusByUid(ctx, pointsLog.Uid, model.StatusQueued)
		if err != nil {
			logger.LogBusinessErr(ctx, errmsg.ErrorDbUpdate, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
			return nil, err
		}
		return &PointsResult{Success: true, Message: "任务已排队"}, nil
	}

	// 发送kafka
	msg := UserPointsMsg{
		Uid:        uid,
		AccountId:  accountId,
		UserId:     userId,
		Amount:     amount,
		RetryTimes: 0,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorJsonMarshal, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
		return nil, err
	}
	if err := svcCtx.KqPusherClient.Push(ctx, string(msgBytes)); err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorKafkaPush, err, logger.WithUserID(fmt.Sprintf("%d", userId)))
		return nil, err
	}

	return &PointsResult{Success: true, Message: "处理中"}, nil
}

// sendDelayCheck 发送延时消息到 dq
func sendDelayCheck(ctx context.Context, svcCtx *svc.ServiceContext, uid int64) error {
	delay := time.Minute * 15

	body := []byte(strconv.FormatInt(uid, 10))

	_, err := svcCtx.DqPusherClient.Delay(body, delay)
	if err != nil {
		logger.LogBusinessErr(ctx, errmsg.ErrorDelayMsg, err, logger.WithUserID(fmt.Sprintf("%d", uid)))
		return err
	}
	return nil
}
