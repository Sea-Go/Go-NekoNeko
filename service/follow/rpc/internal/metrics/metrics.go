/*
 * @Author: Will
 * @Email: haichao.wang@mintegral.com
 * @Date: 2026-02-24 16:05:48
 * @LastEditTime: 2026-02-24 16:05:48
 * @FilePath: /Sea-TryGo/service/follow/rpc/internal/metrics/metrics.go
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	// 路径按你项目实际 config 包调整；函数签名保持和你给的代码一致
	"sea-try-go/service/follow/rpc/internal/config"
)

var (
	// 1) RPC 请求计数：关注系统所有 RPC 方法的调用次数
	// labels:
	// - module: follow_rpc
	// - action: Follow/Unfollow/Block/Unblock/GetFollowList/GetFollowerList/GetBlockList/GetRecommendations
	// - result: ok/biz_fail/sys_fail
	FollowRequestCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "request_total",
		Help:      "",
	}, []string{"module", "action", "result"})

	// 2) RPC 耗时累加（秒）：用 Counter 复照你的写法（更严谨一般用 Histogram，但你要求按原写法）
	// labels:
	// - module: follow_rpc
	// - action: 同上
	FollowRequestSecondsCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "request_seconds_counter",
		Help:      "",
	}, []string{"module", "action"})

	// 3) 关系操作计数：关注/取关/拉黑/解除拉黑
	// labels:
	// - module: follow_relation
	// - action: follow/unfollow/block/unblock
	// - result: ok/fail
	FollowRelationCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "relation_ops_total",
		Help:      "",
	}, []string{"module", "action", "result"})

	// 4) Neo4j 错误计数
	// labels:
	// - module: follow_neo4j
	// - action: FollowUser/UnfollowUser/BlockUser/UnblockUser/GetFollowerList/GetFollowList/GetBlockList/GetRecommendations
	// - type: run/query/scan/parse（你可以按实际分类）
	FollowNeo4jErrorCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "neo4j_error_total",
		Help:      "",
	}, []string{"module", "action", "type"})

	// 5) 列表返回长度（Gauge 存“最近一次”的长度）
	// labels:
	// - module: follow_list
	// - action: following/follower/blocked/recommendation
	FollowListSizeGaugeMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "list_size",
		Help:      "",
	}, []string{"module", "action"})

	// 6) 被拉黑/互拉黑导致的拒绝计数
	// labels:
	// - module: follow_guard
	// - action: Follow/GetRecommendations/...
	FollowBlockedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "follow",
		Name:      "blocked_total",
		Help:      "",
	}, []string{"module", "action"})
)

func InitMetrics(cfg *config.Config) {
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	prometheus.Register(FollowRequestCounterMetric)
	prometheus.Register(FollowRequestSecondsCounterMetric)
	prometheus.Register(FollowRelationCounterMetric)
	prometheus.Register(FollowNeo4jErrorCounterMetric)
	prometheus.Register(FollowListSizeGaugeMetric)
	prometheus.Register(FollowBlockedCounterMetric)
}
