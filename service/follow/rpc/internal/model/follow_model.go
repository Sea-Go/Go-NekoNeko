package model

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Recommendation 推荐结果结构体 (解耦 pb 依赖)
type Recommendation struct {
	TargetId    int64
	MutualScore int64
}

// FollowModel 定义底层数据库所有的 IO 接口
type FollowModel interface {
	FollowUser(ctx context.Context, userId, targetId int64) error
	UnfollowUser(ctx context.Context, userId, targetId int64) error
	BlockUser(ctx context.Context, userId, targetId int64) error
	UnblockUser(ctx context.Context, userId, targetId int64) error
	GetFollowList(ctx context.Context, userId int64, offset, limit int32) ([]int64, error)
	GetFollowerList(ctx context.Context, targetId int64, offset, limit int32) ([]int64, error)
	GetBlockList(ctx context.Context, userId int64, offset, limit int32) ([]int64, error)
	GetRecommendations(ctx context.Context, userId int64, offset, limit int32) ([]*Recommendation, error)
}

type defaultFollowModel struct {
	driver neo4j.DriverWithContext
}

func NewFollowModel(driver neo4j.DriverWithContext) FollowModel {
	return &defaultFollowModel{driver: driver}
}

// 1. 关注操作 (带有属性初始化与亲密度累加)
func (m *defaultFollowModel) FollowUser(ctx context.Context, userId, targetId int64) error {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (u:User {uid: $u_id})
			MERGE (t:User {uid: $t_id})
			WITH u, t
			WHERE NOT (u)-[:BLOCKS]-(t)
			MERGE (u)-[r:FOLLOWS]->(t)
			ON CREATE SET r.created_time = timestamp(), r.weight = 1, r.status = 1
			ON MATCH SET r.weight = r.weight + 1, r.status = 1
			RETURN r
		`
		return tx.Run(ctx, query, map[string]any{
			"u_id": strconv.FormatInt(userId, 10),
			"t_id": strconv.FormatInt(targetId, 10),
		})
	})
	return err
}

// 2. 取消关注 (物理删除边，如果业务需要软删除可改为 SET r.status = 0)
func (m *defaultFollowModel) UnfollowUser(ctx context.Context, userId, targetId int64) error {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {uid: $u_id})-[r:FOLLOWS]->(t:User {uid: $t_id}) DELETE r`
		return tx.Run(ctx, query, map[string]any{
			"u_id": strconv.FormatInt(userId, 10),
			"t_id": strconv.FormatInt(targetId, 10),
		})
	})
	return err
}

// 3. 拉黑操作 (新增黑名单状态，并斩断双向关注)
func (m *defaultFollowModel) BlockUser(ctx context.Context, userId, targetId int64) error {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (u:User {uid: $u_id})
			MERGE (t:User {uid: $t_id})
			MERGE (u)-[b:BLOCKS]->(t)
			ON CREATE SET b.created_time = timestamp(), b.status = 2
			WITH u, t
			MATCH (u)-[f:FOLLOWS]-(t) 
			DELETE f
		`
		return tx.Run(ctx, query, map[string]any{
			"u_id": strconv.FormatInt(userId, 10),
			"t_id": strconv.FormatInt(targetId, 10),
		})
	})
	return err
}

// 4. 取消拉黑
func (m *defaultFollowModel) UnblockUser(ctx context.Context, userId, targetId int64) error {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {uid: $u_id})-[r:BLOCKS]->(t:User {uid: $t_id}) DELETE r`
		return tx.Run(ctx, query, map[string]any{
			"u_id": strconv.FormatInt(userId, 10),
			"t_id": strconv.FormatInt(targetId, 10),
		})
	})
	return err
}

// 5. 获取关注列表 (按亲密度 weight 排序)
func (m *defaultFollowModel) GetFollowList(ctx context.Context, userId int64, offset, limit int32) ([]int64, error) {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {uid: $u_id})-[r:FOLLOWS]->(t:User)
			RETURN t.uid AS target_uid
			ORDER BY r.weight DESC, r.created_time DESC
			SKIP $offset LIMIT $limit
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"u_id":   strconv.FormatInt(userId, 10),
			"offset": offset,
			"limit":  limit,
		})
		if err != nil {
			return nil, err
		}

		var ids []int64
		for result.Next(ctx) {
			uidStr := result.Record().Values[0].(string)
			uid, _ := strconv.ParseInt(uidStr, 10, 64)
			ids = append(ids, uid)
		}
		return ids, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]int64), nil
}

// 6. 获取粉丝列表
func (m *defaultFollowModel) GetFollowerList(ctx context.Context, targetId int64, offset, limit int32) ([]int64, error) {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User)-[r:FOLLOWS]->(t:User {uid: $t_id})
			RETURN u.uid AS user_uid
			ORDER BY r.created_time DESC
			SKIP $offset LIMIT $limit
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"t_id":   strconv.FormatInt(targetId, 10),
			"offset": offset,
			"limit":  limit,
		})
		if err != nil {
			return nil, err
		}

		var ids []int64
		for result.Next(ctx) {
			uidStr := result.Record().Values[0].(string)
			uid, _ := strconv.ParseInt(uidStr, 10, 64)
			ids = append(ids, uid)
		}
		return ids, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]int64), nil
}

// 7. 获取黑名单
func (m *defaultFollowModel) GetBlockList(ctx context.Context, userId int64, offset, limit int32) ([]int64, error) {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {uid: $u_id})-[r:BLOCKS]->(t:User)
			RETURN t.uid AS target_uid
			ORDER BY r.created_time DESC
			SKIP $offset LIMIT $limit
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"u_id":   strconv.FormatInt(userId, 10),
			"offset": offset,
			"limit":  limit,
		})
		if err != nil {
			return nil, err
		}

		var ids []int64
		for result.Next(ctx) {
			uidStr := result.Record().Values[0].(string)
			uid, _ := strconv.ParseInt(uidStr, 10, 64)
			ids = append(ids, uid)
		}
		return ids, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]int64), nil
}

// 8. 推荐算法 (BFS)
func (m *defaultFollowModel) GetRecommendations(ctx context.Context, userId int64, offset, limit int32) ([]*Recommendation, error) {
	session := m.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {uid: $u_id})-[:FOLLOWS*2..3]->(rec:User)
			WHERE u <> rec 
			  AND NOT (u)-[:FOLLOWS]->(rec)
			  AND NOT (u)-[:BLOCKS]-(rec)
			RETURN rec.uid AS target_uid, count(*) AS mutual_score
			ORDER BY mutual_score DESC
			SKIP $offset LIMIT $limit
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"u_id":   strconv.FormatInt(userId, 10),
			"offset": offset,
			"limit":  limit,
		})
		if err != nil {
			return nil, err
		}

		var recs []*Recommendation
		for result.Next(ctx) {
			record := result.Record()
			uidStr := record.Values[0].(string)
			score := record.Values[1].(int64)

			uid, _ := strconv.ParseInt(uidStr, 10, 64)
			recs = append(recs, &Recommendation{TargetId: uid, MutualScore: score})
		}
		return recs, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]*Recommendation), nil
}
