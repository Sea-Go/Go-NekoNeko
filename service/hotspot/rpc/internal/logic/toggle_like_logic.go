package logic

import (
	"context"
	"errors"
	"sea-try-go/service/hotspot/rpc/internal/svc"
	"sea-try-go/service/hotspot/rpc/model"
	"sea-try-go/service/hotspot/rpc/pb"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ToggleLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewToggleLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleLikeLogic {
	return &ToggleLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞/取消点赞
func (l *ToggleLikeLogic) ToggleLike(in *__.LikeReq) (*__.LikeResp, error) {
	// todo: add your logic here and delete this line
	bizID, err := strconv.ParseInt(in.BizId, 10, 64)
	if err != nil {
		return nil, errors.New("invalid biz_id")
	}
	userID, err := strconv.ParseInt(in.UserId, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user_id")
	}

	// 2. 事务处理
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 检查当前是否已点赞
		var existingLike model.Like
		result := tx.Where("biz_id = ? AND user_id = ?", bizID, userID).First(&existingLike)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		isLiked := result.Error == nil

		// 行为判断
		if in.Action == int32(__.TypeSet_Upvote_Like) {
			if isLiked {
				return nil // 已经点赞过, 懒得获取 count 直接返回 null
			}
			// 插入点赞记录
			newLike := model.Like{
				BizID:  bizID,
				UserID: userID,
			}
			if err := tx.Create(&newLike).Error; err != nil {
				return err
			}
		} else if in.Action == int32(__.TypeSet_Downvote_Dislike) {
			if !isLiked { // 没点赞取消赞
				return nil
			}
			// 删除点赞记录
			if err := tx.Delete(&existingLike).Error; err != nil {
				return err
			}
			// 插入取消日志
			cancelLog := model.CancelLikesLog{
				UserID: userID,
				BizID:  bizID,
			}
			if err := tx.Create(&cancelLog).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		l.Logger.Errorf("ToggleLike failed: %v", err)
		return nil, err
	}

	// 统计当前点赞总数 (TODO: 从 Redis 获取)
	var count int64
	l.svcCtx.DB.Model(&model.Like{}).Where("biz_id = ?", bizID).Count(&count)

	return &__.LikeResp{
		LikeCount: count,
	}, nil
}
