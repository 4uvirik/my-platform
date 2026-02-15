package observability

import "github.com/prometheus/client_golang/prometheus"

var KafkaProcessingLatency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: "notification_service",
		Name:      "event_processing_seconds",
		Help:      "Kafka event processing latency",
		Buckets:   prometheus.DefBuckets,
	},
)

func init() {
	prometheus.MustRegister(KafkaProcessingLatency)
}
