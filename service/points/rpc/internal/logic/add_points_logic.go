package logic

import (
	"context"
	"encoding/json"
	"errors"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/points/rpc/internal/model"
	"strconv"
	"time"

	"sea-try-go/service/points/rpc/internal/svc"
	pb "sea-try-go/service/points/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPointsMsg struct {
	Uid        int64 `json:"uid"`
	AccountId  int64 `json:"accountId"`
	UserId     int64 `json:"userId"`
	Amount     int32 `json:"amount"`
	RetryTimes int32 `json:"retryTimes"`
}

type AddPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPointsLogic {
	return &AddPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPointsLogic) AddPoints(in *pb.AddPointsReq) (*pb.AddPointsResp, error) {
	// todo: add your logic here and delete this line
	userId := in.UserId
	accountId := in.RequestId

	// 使用雪花算法生成唯一 Uid
	uid, err := snowflake.GetID()
	if err != nil {
		l.Logger.Errorf("生成雪花ID失败: %v", err)
		return nil, err
	}

	// Step 1: 检查是否有正在处理的任务
	hasProcessing, err := l.svcCtx.PointsModel.HasProcessingByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	if hasProcessing {
		// 有人排队 -> 插入 Queued 记录并返回
		err := l.svcCtx.PointsModel.Insert(l.ctx, &model.Points{
			Uid:       uid,
			UserId:    userId,
			AccountId: accountId,
			Amount:    in.AddPoints,
			Status:    model.StatusQueued,
		})
		if err != nil {
			return nil, err
		}
		return &pb.AddPointsResp{Success: true, Message: "任务已排队"}, nil
	}

	// 无人排队 -> 开始处理流程
	// Step 2: 开启事务
	tx := l.svcCtx.PointsModel.BeginTransaction()
	pointsLog := &model.Points{
		Uid:       uid, // 雪花算法生成的唯一ID
		AccountId: accountId,
		UserId:    userId,
		Amount:    in.AddPoints,
		Status:    model.StatusProcessing,
	}
	if err := tx.Create(pointsLog).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Step 4: 发送延时消息 (带重试)
	for i := 0; i < 3; i++ {
		if err := l.sendDelayCheck(pointsLog.Uid); err == nil {
			break
		}
		if i == 2 {
			tx.Rollback()
			return nil, errors.New("发送延时消息失败")
		}
	}

	// Step 5: 提交事务
	tx.Commit()

	// Step 6: Double Check
	hasOther, _ := l.svcCtx.PointsModel.HasOtherProcessingByUserId(l.ctx, userId, pointsLog.Uid)
	if hasOther {
		// 发现冲突，降级为 Queued
		err := l.svcCtx.PointsModel.UpdateStatusByUid(l.ctx, pointsLog.Uid, model.StatusQueued)
		if err != nil {
			return nil, err
		}
		return &pb.AddPointsResp{Success: true, Message: "任务已排队"}, nil
	}

	// 发送kafka
	msg := UserPointsMsg{
		Uid:        uid,
		AccountId:  accountId,
		UserId:     userId,
		Amount:     in.AddPoints,
		RetryTimes: 0,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(msgBytes)); err != nil {
		return nil, err
	}

	return &pb.AddPointsResp{Success: true, Message: "处理中"}, nil
}

// sendDelayCheck 发送延时消息到 dq
func (l *AddPointsLogic) sendDelayCheck(uid int64) error {
	delay := time.Second * 5

	body := []byte(strconv.FormatInt(uid, 10))

	_, err := l.svcCtx.DqPusherClient.Delay(body, delay)
	if err != nil {
		return err
	}
	return nil
}
