package logic

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// 2. 取消关注 (O(1) 删除边)
func (l *UnfollowLogic) Unfollow(in *pb.RelationReq) (*pb.BaseResp, error) {
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(l.ctx)

	_, err := session.ExecuteWrite(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {uid: $u_id})-[r:FOLLOWS]->(t:User {uid: $t_id}) DELETE r`
		return tx.Run(l.ctx, query, map[string]any{
			"u_id": strconv.FormatInt(in.UserId, 10),
			"t_id": strconv.FormatInt(in.TargetId, 10),
		})
	})
	if err != nil {
		l.Logger.Errorf("Unfollow failed: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "Internal Error"}, err
	}
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
