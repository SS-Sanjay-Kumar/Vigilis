package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vigilis_http_requests_total",
			Help: "Total number of http requests",
		},
		[]string{"method", "path", "status"},
	)

	LogsIngestedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "vigilis_logs_ingested_total",
			Help: "Total raw log entries submitted via API",
		},
	)
)
