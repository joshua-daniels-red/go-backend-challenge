package stream

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	once sync.Once

	EventsConsumedFromStream = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_consumed_from_stream_total",
			Help: "Number of events consumed from the Wikipedia stream",
		},
	)
	EventsProducedToRedpanda = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_produced_to_redpanda_total",
			Help: "Number of events produced to Redpanda",
		},
	)
	EventsConsumedFromRedpanda = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_consumed_from_redpanda_total",
			Help: "Number of events consumed from Redpanda",
		},
	)
	EventsProcessedSuccessfully = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_processed_successfully_total",
			Help: "Number of events processed successfully",
		},
	)
	EventsFailedToProcess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_failed_to_process_total",
			Help: "Number of events that failed during processing",
		},
	)
)

func RegisterMetrics() {
	once.Do(func() {
		prometheus.MustRegister(
			EventsConsumedFromStream,
			EventsProducedToRedpanda,
			EventsConsumedFromRedpanda,
			EventsProcessedSuccessfully,
			EventsFailedToProcess,
		)
	})
}
