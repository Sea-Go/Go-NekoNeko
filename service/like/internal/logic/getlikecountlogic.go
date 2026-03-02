package logic

import (
	"context"
	"fmt"
	"strconv"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeCountLogic {
	return &GetLikeCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLikeCountLogic) GetLikeCount(in *pb.GetLikeCountReq) (*pb.GetLikeCountResp, error) {
	resp := &pb.GetLikeCountResp{
		Counts: make(map[string]*pb.LikeCountItem),
	}
	if len(in.TargetIds) == 0 {
		return resp, nil
	}

	keys := make([]string, 0, len(in.TargetIds)*2)
	for _, id := range in.TargetIds {
		keys = append(keys, fmt.Sprintf("like_cnt:%s:%s", in.TargetType, id))
		keys = append(keys, fmt.Sprintf("dislike_cnt:%s:%s", in.TargetType, id))
	}
	vals, err := l.svcCtx.BizRedis.MgetCtx(l.ctx, keys...)
	//如果缓存里没读到就返回"",数量是永远不会变的
	if err != nil {
		l.Errorf("Redis MGET失败:%v", err)
	}
	var missingIds []string
	for i, id := range in.TargetIds {
		var likeStr, dislikeStr string
		if len(vals) > 2*i+1 {
			likeStr = vals[i*2]
			dislikeStr = vals[i*2+1]
		}
		if likeStr == "" || dislikeStr == "" {
			missingIds = append(missingIds, id)
			continue
		}
		likeCount, _ := strconv.ParseInt(likeStr, 10, 64)
		dislikeCount, _ := strconv.ParseInt(dislikeStr, 10, 64)
		resp.Counts[id] = &pb.LikeCountItem{
			LikeCount:    likeCount,
			DislikeCount: dislikeCount,
		}
	}
	if len(missingIds) > 0 {
		l.Infof("触发DB消息回源,缺失数量:%d", len(missingIds))
		dbResMap, dbErr := l.svcCtx.LikeModel.GetBatchLikeCount(in.TargetType, missingIds)
		if dbErr != nil {
			l.Errorf("DB批量查询失败:%v", dbErr)
			return nil, dbErr
		}
		ttl := l.svcCtx.Config.Storage.Redis.CacheTTL
		for _, id := range missingIds {
			likeCount := dbResMap[id][1]
			dislikeCount := dbResMap[id][2]
			resp.Counts[id] = &pb.LikeCountItem{
				LikeCount:    likeCount,
				DislikeCount: dislikeCount,
			}
			likeKey := fmt.Sprintf("like_cnt:%s:%s", in.TargetType, id)
			dislikeKey := fmt.Sprintf("dislike_cnt:%s:%s", in.TargetType, id)
			//trconv.FormatInt(likeCount, 10)中的10表示十进制
			if err := l.svcCtx.BizRedis.SetexCtx(l.ctx, likeKey, strconv.FormatInt(likeCount, 10), ttl); err != nil {
				l.Errorf("写出Redis出错:%v", err)
			}
			if err := l.svcCtx.BizRedis.SetexCtx(l.ctx, dislikeKey, strconv.FormatInt(dislikeCount, 10), ttl); err != nil {
				l.Errorf("写入Redis出错:%v", err)
			}
		}
	}
	return resp, nil
}
