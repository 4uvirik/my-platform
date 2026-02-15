package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	GRPCRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "order_service",
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		},
		[]string{"method"},
	)

	GRPCLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "order_service",
			Name:      "grpc_latency_seconds",
			Help:      "gRPC request latency",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	GRPCErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "order_service",
			Name:      "grpc_errors_total",
			Help:      "Total number of gRPC errors",
		},
		[]string{"method"},
	)
)

func Register() {
	prometheus.MustRegister(GRPCRequests, GRPCLatency, GRPCErrors)
}
