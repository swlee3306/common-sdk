package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	MessagesSent     prometheus.Counter
	MessagesReceived prometheus.Counter
	MessageSize      prometheus.Histogram
	ProcessingTime   prometheus.Histogram
	ActiveConnections prometheus.Gauge
	Errors           prometheus.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		MessagesSent: promauto.NewCounter(prometheus.CounterOpts{
			Name: "multicast_messages_sent_total",
			Help: "Total number of messages sent",
		}),
		MessagesReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "multicast_messages_received_total",
			Help: "Total number of messages received",
		}),
		MessageSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "multicast_message_size_bytes",
			Help: "Size of multicast messages in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		}),
		ProcessingTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "multicast_processing_duration_seconds",
			Help: "Time spent processing messages",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		}),
		ActiveConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "multicast_active_connections",
			Help: "Number of active multicast connections",
		}),
		Errors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "multicast_errors_total",
			Help: "Total number of multicast errors",
		}),
	}
}

func (m *Metrics) RecordMessageSent(size int) {
	m.MessagesSent.Inc()
	m.MessageSize.Observe(float64(size))
}

func (m *Metrics) RecordMessageReceived(size int) {
	m.MessagesReceived.Inc()
	m.MessageSize.Observe(float64(size))
}

func (m *Metrics) RecordProcessingTime(duration time.Duration) {
	m.ProcessingTime.Observe(duration.Seconds())
}

func (m *Metrics) SetActiveConnections(count int) {
	m.ActiveConnections.Set(float64(count))
}

func (m *Metrics) RecordError() {
	m.Errors.Inc()
}

func (m *Metrics) GetHandler() http.Handler {
	return promhttp.Handler()
}
