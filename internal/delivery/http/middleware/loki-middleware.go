package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

type LokiMiddleware struct {
	lokiClient *loki.Client
	logger     *zap.Logger
}

func NewLokiMiddleware(lokiClient *loki.Client, logger *zap.Logger) *LokiMiddleware {
	return &LokiMiddleware{
		lokiClient: lokiClient,
		logger:     logger,
	}
}

func (m *LokiMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		err := m.lokiClient.Handle(model.LabelSet{}, time.Now(), m.formatLogLine(r, duration))
		if err != nil {
			m.logger.Error("Failed to push log entry to Loki", zap.Error(err))
		}
	})
}

func (m *LokiMiddleware) formatLogLine(r *http.Request, duration time.Duration) string {
	logLine := fmt.Sprintf("Method: %s, URL: %s, Duration: %s", r.Method, r.URL.Path, duration)
	m.logger.Sugar().Info(logLine)
	return logLine
}
