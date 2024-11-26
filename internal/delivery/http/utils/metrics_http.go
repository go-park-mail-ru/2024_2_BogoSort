package utils

// import (
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// )

// type HTTPMetrics struct {
// 	totalHits   *prometheus.CounterVec
// 	serviceName string
// 	duration    *prometheus.HistogramVec
// 	totalErrors *prometheus.CounterVec
// }

// func CreateHTTPMetrics(service string) (*HTTPMetrics, error) {
// 	var metric HTTPMetrics
// 	metric.serviceName = service

// 	metric.totalHits = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: service + "_total_hits_count",
// 			Help: "Number of total HTTP requests",
// 		},
// 		[]string{"path", "service", "code"},
// 	)
// 	if err := prometheus.Register(metric.totalHits); err != nil {
// 		return nil, err
// 	}

// 	metric.duration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    service + "_request_duration_seconds",
// 			Help:    "Duration of HTTP requests in seconds",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"path", "service", "code"},
// 	)
// 	if err := prometheus.Register(metric.duration); err != nil {
// 		return nil, err
// 	}

// 	metric.totalErrors = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: service + "_errors_total",
// 			Help: "Number of total HTTP errors",
// 		},
// 		[]string{"path", "service"},
// 	)
// 	if err := prometheus.Register(metric.totalErrors); err != nil {
// 		return nil, err
// 	}

// 	return &metric, nil
// }

// func (m *HTTPMetrics) IncTotalHits(path, code string) {
// 	m.totalHits.WithLabelValues(path, m.serviceName, code).Inc()
// }

// func (m *HTTPMetrics) AddRequestDuration(path, code string, duration time.Duration) {
// 	m.duration.WithLabelValues(path, m.serviceName, code).Observe(duration.Seconds())
// }

// func (m *HTTPMetrics) IncTotalErrors(path string) {
// 	m.totalErrors.WithLabelValues(path, m.serviceName).Inc()
// }
