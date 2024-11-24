package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type HTTPMetrics struct {
	totalHits   *prometheus.CounterVec
	totalErrors *prometheus.CounterVec
	serviceName string
	duration    *prometheus.HistogramVec
}

func CreateHTTPMetrics(service string) (*HTTPMetrics, error) {
	var metric HTTPMetrics
	metric.serviceName = service

	metric.totalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: service + "_total_hits_count",
			Help: "Number of total http requests",
		},
		[]string{"path", "service", "code"})
	if err := prometheus.Register(metric.totalHits); err != nil {
		return nil, err
	}

	metric.totalErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: service + "_error_hits_count",
			Help: "Number of total http error requests",
		},
		[]string{"path", "service", "code"})
	if err := prometheus.Register(metric.totalErrors); err != nil {
		return nil, err
	}

	metric.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: service + "_code",
			Help: "Request time",
		},
		[]string{"path", "service", "code"})
	if err := prometheus.Register(metric.duration); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m *HTTPMetrics) IncTotalHits(path, code string) {
	m.totalHits.WithLabelValues(path, m.serviceName, code).Inc()
}

func (m *HTTPMetrics) IncTotalErrors(path, code string) {
	m.totalErrors.WithLabelValues(path, m.serviceName, code).Inc()
}

func (m *HTTPMetrics) AddRequestDuration(path, code string, duration time.Duration) {
	m.duration.WithLabelValues(path, m.serviceName, code).Observe(duration.Seconds())
}