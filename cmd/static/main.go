package main

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/interceptors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
	staticProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer func() {
		if err := zap.L().Sync(); err != nil {
			zap.L().Error("Failed to sync logger", zap.Error(err))
		}
	}()

	zap.L().Info("Starting static server")

	cfg, err := config.Init()
	if err != nil {
		zap.L().Error("Failed to init config", zap.Error(err))
	}

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL(), int32(cfg.GetPGMaxConns()))
	if err != nil {
		zap.L().Error("Failed to connect to Postgres", zap.Error(err))
	}

	staticRepo, err := postgres.NewStaticRepository(context.Background(), dbPool, cfg.Static.Path, cfg.Static.MaxSize, zap.L(), cfg.PGTimeout)
	if err != nil {
		zap.L().Error("Failed to create static repository", zap.Error(err))
	}

	staticUseCase := service.NewStaticService(staticRepo)

	metrics, err := metrics.NewGRPCMetrics("static")
	if err != nil {
		zap.L().Fatal("Ошибка при инициализации метрик", zap.Error(err))
	}

	staticService := static.NewStaticGrpc(staticUseCase)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewMetricsInterceptor(*metrics).NewMetricsInterceptor),
	)
	staticProto.RegisterStaticServiceServer(server, staticService)
	addr := cfg.StaticHost + ":" + strconv.Itoa(cfg.StaticPort)

	http.Handle("/api/v1/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":7053", nil); err != nil {
			zap.L().Fatal("Failed to start metrics HTTP server", zap.Error(err))
		}
	}()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Error("Failed to listen on port", zap.Error(err))
	}
	zap.L().Info("Listening on grpc address", zap.String("address", addr))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	go func() {
		err := server.Serve(lis)
		if err != nil {
			zap.L().Error("Failed to serve grpc", zap.Error(err))
		}
	}()
	<-ctx.Done()
	server.GracefulStop()
}
