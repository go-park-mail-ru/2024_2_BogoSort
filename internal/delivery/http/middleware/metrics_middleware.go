package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/gorilla/mux"
)

func CreateMetricsMiddleware(metric *metrics.HTTPMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			rec := &statusRecorder{ResponseWriter: writer, statusCode: 200}

			start := time.Now()

			next.ServeHTTP(rec, request)

			end := time.Since(start)

			codeStr := strconv.Itoa(rec.statusCode)
			route := mux.CurrentRoute(request)
			path, _ := route.GetPathTemplate()
			method := request.Method

			if path != "/api/v1/metrics" {
				metric.AddRequestDuration(path, method, codeStr, end)
				metric.IncTotalHits(path, method, codeStr)
				if rec.statusCode >= 400 {
					metric.IncTotalErrors(path, method, codeStr)
				}
			}
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
