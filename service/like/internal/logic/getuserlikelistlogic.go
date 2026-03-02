package logic

import (
	"context"
	"fmt"
	"strconv"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLikeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeListLogic {
	return &GetUserLikeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLikeListLogic) GetUserLikeList(in *pb.GetUserLikeListReq) (*pb.GetUserLikeListResp, error) {

	limit := int(in.PageSize)
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	redisKey := fmt.Sprintf("user_like_list:%d:%s", in.UserId, in.TargetType)
	maxScore := "+inf"
	if in.Cursor > 0 {
		//"("表示不包含当前游标,相当于SQL中的 <
		maxScore = fmt.Sprintf("(%d", in.Cursor)
	}
	script := `return redis.call("ZREVRANGEBYSCORE", KEYS[1], ARGV[1], "-inf", "WITHSCORES", "LIMIT", "0", ARGV[2])`
	res, err := l.svcCtx.BizRedis.EvalCtx(l.ctx, script, []string{redisKey}, maxScore, strconv.Itoa(limit))
	if err != nil {
		l.Errorf("查询失败:%v", err)
	}
	var list []*pb.LikeRecordItem
	var nextCursor int64
	if res != nil {
		if vals, ok := res.([]interface{}); ok {
			for i := 0; i < len(vals); i += 2 {
				var targetId string
				var score int64

				switch val := vals[i].(type) {
				case string:
					targetId = val
				case []byte:
					targetId = string(val)
				}
				switch val := vals[i+1].(type) {
				case string:
					score, _ = strconv.ParseInt(val, 10, 64)
				case []byte:
					score, _ = strconv.ParseInt(string(val), 10, 64)
				}
				realTimestamp := (score >> 22)
				list = append(list, &pb.LikeRecordItem{
					TargetId:  targetId,
					Timestamp: realTimestamp,
				})
				nextCursor = score
			}
		}
	}
	if len(list) == 0 {
		l.Infof("触发DB回源")
		dbResults, dbErr := l.svcCtx.LikeModel.GetUserLikeList(in.UserId, in.TargetType, in.Cursor, limit)
		if dbErr != nil {
			return nil, dbErr
		}
		for _, r := range dbResults {
			var timestamp int64 = r.CreateTime
			if timestamp == 0 {
				return nil, fmt.Errorf("系统数据异常")
			}
			list = append(list, &pb.LikeRecordItem{
				TargetId:  r.TargetId,
				Timestamp: r.CreateTime,
			})
			nextCursor = r.Id
		}
	}
	isEnd := len(list) < limit
	return &pb.GetUserLikeListResp{
		List:       list,
		IsEnd:      isEnd,
		NextCursor: nextCursor,
	}, nil
}
