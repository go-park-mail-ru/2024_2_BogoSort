package interceptors

import (
	"context"
	"net/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

type Interceptor struct {
	metrics metrics.GRPCMetrics
}

func CreateMetricsInterceptor(metrics metrics.GRPCMetrics) *Interceptor {
	return &Interceptor{
		metrics: metrics,
	}
}

func getCode(err string) int {
	switch err {
	case "session does not exist":
		return http.StatusUnauthorized
	case "no such cookie in userStorage":
		return http.StatusUnauthorized
	case "user not authorized":
		return http.StatusUnauthorized
	}

	return http.StatusBadRequest
}

func (interceptor *Interceptor) ServeMetricsInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	end := time.Since(start)
	code := http.StatusOK

	if err != nil {
		code = getCode(err.Error())
	}

	codeStr := strconv.Itoa(code)
	interceptor.metrics.AddDuration(codeStr, info.FullMethod, end)
	interceptor.metrics.IncTotalHits(codeStr, info.FullMethod)
	if code >= 400 {
		interceptor.metrics.IncTotalErrors(codeStr, info.FullMethod)
	}

	return h, err
}