package interceptors

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"google.golang.org/grpc"
	"go.uber.org/zap"
	"strconv"
)

type MetricsInterceptor struct {
	metrics metrics.GRPCMetrics
}

func NewMetricsInterceptor(metrics metrics.GRPCMetrics) *MetricsInterceptor {
	return &MetricsInterceptor{metrics: metrics}
}

func (m *MetricsInterceptor) ServeMetricsClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Since(start)

	code := http.StatusOK
	if err != nil {
		code = getCode(err.Error())
	}

	zap.L().Info("method", zap.String("method", method))
	codeStr := strconv.Itoa(code)
	m.metrics.AddDuration(codeStr, method, duration)
	m.metrics.IncTotalHits(codeStr, method)
	if code >= 400 {
		m.metrics.IncTotalErrors(codeStr, method)
	}

	return err
}

func getCode(err string) int {
	switch err {
	case "session does not exist":
		return http.StatusUnauthorized
	case "no such cookie in userStorage":
		return http.StatusUnauthorized
	case "user not authorized":
		return http.StatusUnauthorized
	case "invalid input":
		return http.StatusBadRequest
	case "resource not found":
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}