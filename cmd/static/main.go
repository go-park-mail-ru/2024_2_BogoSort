package main

import (
	"context"
	"go.uber.org/zap"
	"net"
	"os/signal"
	"syscall"
	"strconv"
	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"google.golang.org/grpc"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
)

func main() {
	logger := zap.L()
	defer logger.Sync()

	logger.Info("Starting static server")

	cfg, err := config.Init()
	if err != nil {
		logger.Error("Failed to init config", zap.Error(err))
	}

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		logger.Error("Failed to connect to Postgres", zap.Error(err))
	}

	staticRepo, err := postgres.NewStaticRepository(context.Background(), dbPool, cfg.Static.Path, cfg.Static.MaxSize, zap.L(), cfg.PGTimeout)
	if err != nil {
		logger.Error("Failed to create static repository", zap.Error(err))
	}

	staticUseCase := service.NewStaticService(staticRepo, zap.L())
	staticService := static.NewStaticGrpc(staticUseCase)
	server := grpc.NewServer()
	staticProto.RegisterStaticServiceServer(server, staticService)
	addr := cfg.StaticHost + ":" + strconv.Itoa(cfg.StaticPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Failed to listen on port", zap.Error(err))
	}
	logger.Info("Listening on grpc address", zap.String("address", addr))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	go func() {
		err := server.Serve(lis)
		if err != nil {
			logger.Error("Failed to serve grpc", zap.Error(err))
		}
	}()
	<-ctx.Done()
	server.GracefulStop()
}