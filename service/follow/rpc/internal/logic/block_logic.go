package logic

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockLogic {
	return &BlockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Block 3. 拉入黑名单 (O(1) 插入边并清理双向 Follow 边)
func (l *BlockLogic) Block(in *pb.RelationReq) (*pb.BaseResp, error) {
	// Corner Case 特判：不能自己拉黑自己
	if in.UserId == in.TargetId {
		return &pb.BaseResp{Code: 400, Msg: "Cannot block yourself"}, nil
	}

	// 开启 Neo4j 写入会话 (防 MLE/连接泄露，必须 defer 释放)
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(l.ctx)

	// 执行事务
	_, err := session.ExecuteWrite(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Cypher 逻辑 (防 WA 核心)：
		// 1. MERGE 创建 BLOCKS 有向边
		// 2. 匹配两人之间可能存在的双向 FOLLOWS 边，强制 DELETE 斩断
		query := `
			MERGE (u:User {uid: $u_id})
			MERGE (t:User {uid: $t_id})
			MERGE (u)-[b:BLOCKS]->(t)
			WITH u, t
			MATCH (u)-[f:FOLLOWS]-(t) 
			DELETE f
		`
		params := map[string]any{
			"u_id": strconv.FormatInt(in.UserId, 10),
			"t_id": strconv.FormatInt(in.TargetId, 10),
		}

		result, err := tx.Run(l.ctx, query, params)
		if err != nil {
			return nil, err
		}

		return result.Consume(l.ctx) // 确保数据完全落盘
	})

	if err != nil {
		l.Logger.Errorf("Failed to execute Block cypher: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "Internal server error"}, err
	}

	// 状态转移成功 AC
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
