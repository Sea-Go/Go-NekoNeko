package logic

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strconv"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlockListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockListLogic {
	return &GetBlockListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// 7. 获取黑名单列表 (O(K), 遍历 BLOCKS 出边)
func (l *GetBlockListLogic) GetBlockList(in *pb.ListReq) (*pb.UserListResp, error) {
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(l.ctx)

	res, err := session.ExecuteRead(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {uid: $u_id})-[:BLOCKS]->(t:User)
			RETURN t.uid AS target_uid
			SKIP $offset LIMIT $limit
		`
		params := map[string]any{
			"u_id":   fmt.Sprintf("%d", in.UserId),
			"offset": in.Offset,
			"limit":  in.Limit,
		}

		result, err := tx.Run(l.ctx, query, params)
		if err != nil {
			return nil, err
		}

		var userIds []int64
		for result.Next(l.ctx) {
			record := result.Record()
			uidStr := record.Values[0].(string)
			uid, _ := strconv.ParseInt(uidStr, 10, 64)
			userIds = append(userIds, uid)
		}
		return userIds, result.Err()
	})

	if err != nil {
		return &pb.UserListResp{Code: 500, Msg: "DB Error"}, err
	}

	return &pb.UserListResp{Code: 0, Msg: "success", UserIds: res.([]int64)}, nil
}
