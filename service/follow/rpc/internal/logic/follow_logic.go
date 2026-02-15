package logic

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Follow 1. 点击关注
func (l *FollowLogic) Follow(in *pb.RelationReq) (*pb.BaseResp, error) {
	// Corner Case 特判：不能关注自己
	if in.UserId == in.TargetId {
		return &pb.BaseResp{Code: 400, Msg: "User cannot follow themselves"}, nil
	}

	// 开启 Neo4j 写入会话 (防 MLE/连接泄露，必须 defer 释放)
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(l.ctx)

	// 执行事务
	_, err := session.ExecuteWrite(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Cypher 逻辑：
		// 1. MATCH 确保两人没有 BLOCKS 关系 (黑名单拦截)
		// 2. MERGE 确保边是唯一的 (幂等性保障)
		query := `
			MERGE (u:User {uid: $u_id})
			MERGE (t:User {uid: $t_id})
			WITH u, t
			WHERE NOT (u)-[:BLOCKS]-(t)
			MERGE (u)-[r:FOLLOWS]->(t)
			RETURN r
		`
		params := map[string]any{
			"u_id": strconv.FormatInt(in.UserId, 10), // 将 int64 转为 string，保持底层统一
			"t_id": strconv.FormatInt(in.TargetId, 10),
		}

		result, err := tx.Run(l.ctx, query, params)
		if err != nil {
			return nil, err
		}

		return result.Consume(l.ctx) // 确保数据落盘
	})

	if err != nil {
		l.Logger.Errorf("Failed to execute Follow cypher: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "Internal server error"}, err
	}

	// 状态转移成功 AC
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
