package main

import (
	"context"
	"net/http"
	_ "net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/auth"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase"
	static "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
	http3 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/gorilla/mux"
	_ "github.com/grafana/loki-client-go/loki"
	_ "github.com/grafana/loki-client-go/pkg/urlutil"
	"github.com/microcosm-cc/bluemonday"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	// lokiURL, err := url.Parse("http://loki:3100/loki/api/v1/push")
	// if err != nil {
	// 	logger.Fatal("failed to parse loki url", zap.Error(err))
	// }
	// lokiConfig := loki.Config{
	// 	URL: urlutil.URLValue{URL: lokiURL},
	// }
	// lokiClient, err := loki.New(lokiConfig)
	// if err != nil {
	// 	logger.Fatal("failed to create loki client", zap.Error(err))
	// }

	cfg, err := config.Init()
	if err != nil {
		logger.Error("failed to init config", zap.Error(err))
	}

	router, err := Init(cfg)
	if err != nil {
		logger.Error("failed to initialize router", zap.Error(err))
	}

	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggerMiddleware)
	// router.Use(middleware.NewLokiMiddleware(lokiClient, logger).Handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://5.188.141.136:8008",
			"http://localhost:8008",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"X-Authenticated", "X-CSRF-Token"},
		AllowCredentials: true,
	}).Handler(router)

	logger.Info("Server started on " + config.GetServerAddress())

	server := &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      corsHandler,
		ReadTimeout:  config.GetReadTimeout(),
		WriteTimeout: config.GetWriteTimeout(),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", zap.Error(err))
		}
	}()

	<-stop
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownTimeout())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exiting")
}

func Init(cfg config.Config) (*mux.Router, error) {
	var logger = zap.L()

	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	metric, err := metrics.NewHTTPMetrics("app")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http metrics")
	}
	metricsMiddleware := middleware.CreateMetricsMiddleware(metric)
	router.Use(metricsMiddleware)

	policy := bluemonday.UGCPolicy()

	authRouter := router.PathPrefix("").Subrouter()

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL(), int32(cfg.GetPGMaxConns()))
	if err != nil {
		zap.L().Error("Failed to connect to Postgres", zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to Postgres")
	}
	rdb, err := connector.GetRedisConnector(cfg.RdAddr, cfg.RdPass, cfg.RdDB)
	if err != nil {
		zap.L().Error("Failed to connect to Redis", zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	ctx := context.Background()

	advertsRepo, err := postgres.NewAdvertRepository(dbPool, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create advert repository")
	}
	categoryRepo, err := postgres.NewCategoryRepository(dbPool, zap.L(), ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create category repository")
	}
	sessionRepo, err := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), ctx, zap.L())
	if err != nil {
		return nil, handleRepoError(err, "unable to create session repository")
	}
	userRepo, err := postgres.NewUserRepository(dbPool, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create user repository")
	}
	sellerRepo, err := postgres.NewSellerRepository(dbPool, ctx)
	if err != nil {
		return nil, handleRepoError(err, "unable to create seller repository")
	}
	historyRepo, err := postgres.NewHistoryRepository(dbPool, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create history repository")
	}
	csrfToken, err := utils.NewAesCryptHashToken(zap.L())
	if err != nil {
		return nil, handleRepoError(err, "unable to create csrf token")
	}
	authGrpcClient, err := auth.NewGrpcClient(config.GetAuthAddress())
	if err != nil {
		return nil, handleRepoError(err, "unable to create grpc client")
	}

	cartPurchaseClient, err := cart_purchase.NewCartPurchaseClient(config.GetCartPurchaseAddress())
	if err != nil {
		return nil, handleRepoError(err, "unable to create cart purchase client")
	}
	staticClient, err := static.NewStaticGrpcClient(config.GetStaticAddress(), cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create static client")
	}

	advertsUseCase := service.NewAdvertService(advertsRepo, sellerRepo, userRepo, historyRepo)
	categoryUseCase := service.NewCategoryService(categoryRepo)
	userUC := service.NewUserService(userRepo, sellerRepo)
	sessionUC := service.NewAuthService(sessionRepo)
	sessionManager := utils.NewSessionManager(authGrpcClient, int(cfg.Session.ExpirationTime.Seconds()), cfg.Session.SecureCookie, logger)
	router.Use(middleware.NewAuthMiddleware(sessionManager).AuthMiddleware)

	advertsHandler := http3.NewAdvertEndpoint(advertsUseCase, *staticClient, sessionManager, policy)
	authHandler := http3.NewAuthEndpoint(sessionUC, sessionManager)
	userHandler := http3.NewUserEndpoint(userUC, sessionUC, sessionManager, *staticClient, policy)
	sellerHandler := http3.NewSellerEndpoint(sellerRepo)
	purchaseHandler := http3.NewPurchaseEndpoint(cartPurchaseClient)
	cartHandler := http3.NewCartEndpoint(cartPurchaseClient)
	categoryHandler := http3.NewCategoryEndpoint(categoryUseCase)
	staticHandler := http3.NewStaticEndpoint(*staticClient)
	historyHandler := http3.NewHistoryEndpoint(historyRepo)

	csrfEndpoints := http3.NewCSRFEndpoint(csrfToken, sessionManager)
	csrfEndpoints.Configure(router)
	userHandler.ConfigureUnprotectedRoutes(router)
	advertsHandler.ConfigureRoutes(router)

	authRouter.Use(middleware.CSRFMiddleware(csrfToken, sessionManager))

	advertsHandler.ConfigureProtectedRoutes(authRouter)
	categoryHandler.ConfigureRoutes(authRouter)
	authHandler.Configure(authRouter)
	userHandler.ConfigureProtectedRoutes(authRouter)
	sellerHandler.Configure(authRouter)
	cartHandler.Configure(authRouter)
	purchaseHandler.ConfigureRoutes(authRouter)
	staticHandler.ConfigureRoutes(router)
	historyHandler.ConfigureRoutes(authRouter)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/api/v1/metrics").Handler(promhttp.Handler())

	return router, nil
}

func handleRepoError(err error, message string) error {
	zap.L().Error(message, zap.Error(err))
	return errors.Wrap(err, message)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("Panic occurred", zap.Any("error", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
