package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func CreateMetricsMiddleware(metric *metrics.HTTPMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			code := new(int)
			*code = 200
			request = request.WithContext(context.WithValue(request.Context(), "code", code))

			start := time.Now()

			next.ServeHTTP(writer, request)

			end := time.Since(start)

			codeStr := strconv.Itoa(*code)
			route := mux.CurrentRoute(request)
			path, _ := route.GetPathTemplate()
			method := request.Method

			if path != "/metrics" {
				metric.AddRequestDuration(path, method, codeStr, end)
				metric.IncTotalHits(path, method, codeStr)
				if *code >= 400 {
					metric.IncTotalErrors(path, method, codeStr)
				}
			}
		})
	}
}