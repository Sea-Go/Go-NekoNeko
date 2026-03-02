package logic

import (
	"context"
	"fmt"
	"strconv"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLikeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeStateLogic {
	return &GetLikeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLikeStateLogic) GetLikeState(in *pb.GetLikeStateReq) (*pb.GetLikeStateResp, error) {
	resp := &pb.GetLikeStateResp{
		States: make(map[string]int32),
	}
	if len(in.TargetIds) == 0 {
		return resp, nil
	}
	keys := make([]string, 0, len(in.TargetIds))
	for _, id := range in.TargetIds {
		keys = append(keys, fmt.Sprintf("like_state:%d:%s:%s", in.UserId, in.TargetType, id))
	}

	vals, err := l.svcCtx.BizRedis.MgetCtx(l.ctx, keys...)
	if err != nil {
		l.Errorf("Redis MGET失败:%v", err)
	}

	var missingIds []string
	for i, id := range in.TargetIds {
		var stateStr string
		if len(vals) > i {
			stateStr = vals[i]
		}
		if stateStr == "" {
			missingIds = append(missingIds, id)
			continue
		}
		state, _ := strconv.ParseInt(stateStr, 10, 32) //32表示最多32位即int32
		resp.States[id] = int32(state)
	}

	if len(missingIds) > 0 {
		l.Infof("触发DB回源,缺失的数量为:%v", len(missingIds))
		dbRespMap, dbErr := l.svcCtx.LikeModel.GetUserBatchLikeState(in.UserId, in.TargetType, missingIds)
		if dbErr != nil {
			l.Errorf("批量查询失败:%v", dbErr)
			return nil, dbErr
		}
		ttl := l.svcCtx.Config.Storage.Redis.CacheTTL
		//防止内存穿透
		for _, id := range missingIds {
			//Go中如果dbRespMap没有找到id对应的值就会返回默认的零值,然后将零值存入Redis中达成防止内存穿透的目的
			//否则如果多次连续发Redis中不存在的数据,一直调用SQL,会导致很大的问题
			state := dbRespMap[id]
			resp.States[id] = state
			redisKey := fmt.Sprintf("like_state:%d:%s:%s", in.UserId, in.TargetType, id)
			//Itoa也是将int类型的state转换为string类型,面对较小类型的int直接用这个
			if err := l.svcCtx.BizRedis.SetexCtx(l.ctx, redisKey, strconv.Itoa(int(state)), ttl); err != nil {
				l.Errorf("写入Redis出错:%v", err)
			}
		}

	}
	return resp, nil
}
