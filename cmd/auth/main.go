package main

import (
	"context"
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	authProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer logger.Sync()

	// Инициализация конфигурации
	cfg, err := config.Init()
	if err != nil {
		logger.Fatal("Ошибка при инициализации конфигурации", zap.Error(err))
	}

	// Подключение к Redis
	rdb, err := connector.GetRedisConnector(cfg.RdAddr, cfg.RdPass, cfg.RdDB)
	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis", zap.Error(err))
	}

	// Инициализация репозитория сессий
	sessionRepo, err := redis.NewSessionRepository(rdb, 10, context.Background(), logger)
	if err != nil {
		logger.Fatal("Ошибка при инициализации репозитория сессий", zap.Error(err))
	}

	// Инициализация сервисов
	authService := service.NewAuthService(sessionRepo, logger)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	authServer := auth.NewGrpsServer(authService)

	// Регистрация AuthService сервера
	authProto.RegisterAuthServiceServer(grpcServer, authServer)

	// Настройка прослушивания порта
	address := ":50051" // Измените на нужный порт при необходимости
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("Не удалось прослушивать порт", zap.Error(err))
	}

	logger.Info("Auth сервер запущен на " + address)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Ошибка при запуске gRPC сервера", zap.Error(err))
	}
}
