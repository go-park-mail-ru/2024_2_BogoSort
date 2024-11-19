package main

import (
	"context"
	"net"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthProto "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := zap.L()
	defer logger.Sync()

	cfg, err := config.Init()
	if err != nil {
		logger.Fatal("Ошибка при инициализации конфигурации", zap.Error(err))
	}

	rdb, err := connector.GetRedisConnector(cfg.RdAddr, cfg.RdPass, cfg.RdDB)
	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis", zap.Error(err))
	}

	sessionRepo, err := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), context.Background(), zap.L())
	if err != nil {
		logger.Fatal("Ошибка при инициализации репозитория сессий", zap.Error(err))
	}

	authService := service.NewAuthService(sessionRepo, zap.L())

	server := grpc.NewServer()
	authServer := auth.NewGrpcServer(authService)

	healthServer := health.NewServer()
	healthProto.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthProto.HealthCheckResponse_SERVING)

	authProto.RegisterAuthServiceServer(server, authServer)

	address := config.GetAuthAddress()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("Не удалось прослушивать порт", zap.Error(err))
	}

	logger.Info("Auth сервер запущен на " + address)
	if err := server.Serve(lis); err != nil {
		logger.Fatal("Ошибка при запуске gRPC сервера", zap.Error(err))
	}
}
