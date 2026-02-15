package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaEventsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "notification_service",
			Name:      "events_processed_total",
			Help:      "Total processed Kafka events",
		},
	)

	KafkaDuplicates = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "notification_service",
			Name:      "events_duplicate_total",
			Help:      "Total duplicate Kafka events",
		},
	)

	KafkaErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "notification_service",
			Name:      "events_errors_total",
			Help:      "Total Kafka processing errors",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		KafkaEventsProcessed,
		KafkaDuplicates,
		KafkaErrors,
	)
}
