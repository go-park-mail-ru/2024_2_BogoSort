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
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/interceptors"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	zap.L().Info("Starting static server")

	cfg, err := config.Init()
	if err != nil {
		zap.L().Error("Failed to init config", zap.Error(err))
	}

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		zap.L().Error("Failed to connect to Postgres", zap.Error(err))
	}

	staticRepo, err := postgres.NewStaticRepository(context.Background(), dbPool, cfg.Static.Path, cfg.Static.MaxSize, zap.L(), cfg.PGTimeout)
	if err != nil {
		zap.L().Error("Failed to create static repository", zap.Error(err))
	}

	staticUseCase := service.NewStaticService(staticRepo, zap.L())

	metrics, err := metrics.NewGRPCMetrics("static")
	if err != nil {
		zap.L().Fatal("Error initializing metrics", zap.Error(err))
	}

	metricsInterceptor := interceptors.NewMetricsInterceptor(*metrics)

	grpcConn, err := grpc.NewClient(
		config.GetAuthAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(metricsInterceptor.ServeMetricsClientInterceptor),
	)
	if err != nil {
		zap.L().Fatal("Error occurred while starting grpc connection on auth service", zap.Error(err))

		return
	}

	defer grpcConn.Close()
	staticService := static.NewStaticGrpc(staticUseCase)

	server := grpc.NewServer()
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
