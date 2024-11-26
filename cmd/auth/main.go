package main

import (
	"context"
	"net"

	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/interceptors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthProto "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	cfg, err := config.Init()
	if err != nil {
		logger.Error("Error initializing configuration", zap.Error(err))
	}

	rdb, err := connector.GetRedisConnector(cfg.RdAddr, cfg.RdPass, cfg.RdDB)
	if err != nil {
		logger.Error("Error connecting to Redis", zap.Error(err))
	}

	sessionRepo, err := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), context.Background(), zap.L())
	if err != nil {
		logger.Error("Error initializing session repository", zap.Error(err))
	}

	authService := service.NewAuthService(sessionRepo)

	metrics, err := metrics.NewGRPCMetrics("auth")
	if err != nil {
		logger.Error("Error initializing metrics", zap.Error(err))
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.CreateMetricsInterceptor(*metrics).ServeMetricsInterceptor),
	)
	authServer := auth.NewGrpcServer(authService)

	healthServer := health.NewServer()
	healthProto.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthProto.HealthCheckResponse_SERVING)

	authProto.RegisterAuthServiceServer(server, authServer)

	http.Handle("/api/v1/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":7051", nil); err != nil {
			logger.Error("Failed to start metrics HTTP server", zap.Error(err))
		}
	}()

	address := config.GetAuthAddress()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on port", zap.Error(err))
	}

	logger.Info("Auth server started on " + address)
	if err := server.Serve(lis); err != nil {
		logger.Error("Error starting gRPC server", zap.Error(err))
	}
}
