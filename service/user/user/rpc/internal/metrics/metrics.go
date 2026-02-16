package metrics

import "github.com/zeromicro/go-zero/core/metric"

// 演示metrics的使用
var RegisterSuccessTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "user",
	Subsystem: "register",
	Name:      "success_total",
	Help:      "successful user registrations total",
	Labels:    []string{}, // 不需要维度就留空
})

var GetUserLatencyMs = metric.NewHistogramVec(&metric.HistogramVecOpts{
	Namespace: "user",
	Subsystem: "get_user",
	Name:      "latency_ms",
	Help:      "get user rpc latency in ms",
	Labels:    []string{"success"},
})
