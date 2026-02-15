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

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// 5. 获取关注列表 (O(K), 遍历出边)
func (l *GetFollowListLogic) GetFollowList(in *pb.ListReq) (*pb.UserListResp, error) {
	// 使用 Read Session 提高并发性能
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(l.ctx)

	res, err := session.ExecuteRead(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// 分页遍历出边
		query := `
			MATCH (u:User {uid: $u_id})-[:FOLLOWS]->(t:User)
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
			uidStr := record.Values[0].(string) // 解析返回的字符串类型
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
