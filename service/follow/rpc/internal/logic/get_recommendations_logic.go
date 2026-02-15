package logic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecommendationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendationsLogic {
	return &GetRecommendationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// 8. 遍历三层关系查找可能喜欢 (受限 BFS)
func (l *GetRecommendationsLogic) GetRecommendations(in *pb.ListReq) (*pb.RecommendResp, error) {
	session := l.svcCtx.Neo4jDriver.NewSession(l.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(l.ctx)

	res, err := session.ExecuteRead(l.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {uid: $u_id})-[:FOLLOWS*2..3]->(rec:User)
			WHERE u <> rec 
			  AND NOT (u)-[:FOLLOWS]->(rec)
			  AND NOT (u)-[:BLOCKS]-(rec)
			RETURN rec.uid AS target_uid, count(*) AS mutual_score
			ORDER BY mutual_score DESC
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

		var recs []*pb.RecommendResp_RecommendUser
		for result.Next(l.ctx) {
			record := result.Record()
			uidStr := record.Values[0].(string)
			score := record.Values[1].(int64) // count(*) 返回的一定是 int64

			uid, _ := strconv.ParseInt(uidStr, 10, 64)

			recs = append(recs, &pb.RecommendResp_RecommendUser{
				TargetId:    uid,
				MutualScore: int32(score),
			})
		}
		return recs, result.Err()
	})

	if err != nil {
		l.Logger.Errorf("Failed to run BFS recommendation: %v", err)
		return &pb.RecommendResp{Code: 500, Msg: "Recommendation logic failed"}, err
	}

	return &pb.RecommendResp{Code: 0, Msg: "success", Users: res.([]*pb.RecommendResp_RecommendUser)}, nil
}
