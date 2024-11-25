package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/loki/pkg/logproto"
	"go.uber.org/zap"
)

type LokiMiddleware struct {
	lokiClient logproto.PusherClient
	logger     *zap.Logger
}

func NewLokiMiddleware(lokiClient logproto.PusherClient, logger *zap.Logger) *LokiMiddleware {
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
		entry := logproto.Entry{
			Timestamp: time.Now(),
			Line:      m.formatLogLine(r, duration),
		}

		pushRequest := &logproto.PushRequest{
			Streams: []logproto.Stream{
				{
					Entries: []logproto.Entry{entry},
				},
			},
		}

		_, err := m.lokiClient.Push(context.Background(), pushRequest)
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
