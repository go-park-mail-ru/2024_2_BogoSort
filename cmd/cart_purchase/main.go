package main

import (
	"context"
	"net"
	"time"

	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase"
	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/interceptors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthProto "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	cfg, err := config.Init()
	if err != nil {
		zap.L().Fatal("Ошибка при инициализации конфигурации", zap.Error(err))
	}

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		zap.L().Error("Failed to connect to Postgres", zap.Error(err))
		return
	}

	cartRepo, err := postgres.NewCartRepository(dbPool, context.Background())
	if err != nil {
		zap.L().Error("Failed to create cart repository", zap.Error(err))
		return
	}
	advertRepo, err := postgres.NewAdvertRepository(dbPool, context.Background(), time.Duration(cfg.PGTimeout))
	if err != nil {
		zap.L().Error("Failed to create advert repository", zap.Error(err))
		return
	}
	purchaseRepo, err := postgres.NewPurchaseRepository(dbPool, context.Background(), time.Duration(cfg.PGTimeout))
	if err != nil {
		zap.L().Error("Failed to create purchase repository", zap.Error(err))
		return
	}

	metrics, err := metrics.NewGRPCMetrics("cart_purchase")
	if err != nil {
		zap.L().Fatal("Ошибка при инициализации метрик", zap.Error(err))
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

	server := grpc.NewServer()
	cartUC := service.NewCartService(cartRepo, advertRepo, zap.L())
	purchaseUC := service.NewPurchaseService(purchaseRepo, advertRepo, cartRepo, zap.L())
	cartPurchaseServer := cart_purchase.NewGrpcServer(cartUC, purchaseUC)

	healthServer := health.NewServer()
	healthProto.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthProto.HealthCheckResponse_SERVING)

	cartPurchaseProto.RegisterCartPurchaseServiceServer(server, cartPurchaseServer)
	address := config.GetCartPurchaseAddress()

	http.Handle("/api/v1/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":7052", nil); err != nil {
			zap.L().Fatal("Failed to start metrics HTTP server", zap.Error(err))
		}
	}()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("Failed to listen on port", zap.Error(err))
	}

	zap.L().Info("CartPurchase server started on " + address)

	if err := server.Serve(lis); err != nil {
		zap.L().Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
