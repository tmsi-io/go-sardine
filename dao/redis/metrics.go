package redis

import (
	"github.com/tmsi-io/go-sardine/metric"
)

const MetricNameSpace = "sardine_redis_client"

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: MetricNameSpace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "redis client requests duration(ms).",
		Labels:    []string{"addr", "command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: MetricNameSpace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "redis client requests error count.",
		Labels:    []string{"addr", "command", "error"},
	})
	_metricHits = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: MetricNameSpace,
		Subsystem: "",
		Name:      "hits_total",
		Help:      "redis client hits total.",
		Labels:    []string{"addr"},
	})
	_metricMisses = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: MetricNameSpace,
		Subsystem: "",
		Name:      "misses_total",
		Help:      "redis client misses total.",
		Labels:    []string{"addr"},
	})
)
