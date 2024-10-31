package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	http3 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

func NewRouter(cfg config.Config) (*mux.Router, error) {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
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
	userRepo := postgres.NewUserRepository(dbPool, ctx, zap.L())
	sellerRepo := postgres.NewSellerRepository(dbPool, ctx, zap.L())
	userUC := service.NewUserService(userRepo, sellerRepo, zap.L())
	sessionRepo := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), zap.L())
	sessionUC := service.NewAuthService(sessionRepo, zap.L())

	sessionManager := utils.NewSessionManager(sessionUC, int(cfg.Session.ExpirationTime.Seconds()), cfg.Session.SecureCookie, zap.L())

	authHandler := http3.NewAuthEndpoints(sessionUC, sessionManager, zap.L())
	userHandler := http3.NewUserEndpoints(userUC, sessionUC, sessionManager, zap.L())
	sellerHandler := http3.NewSellerEndpoints(sellerRepo, zap.L())

	authHandler.Configure(router)
	userHandler.Configure(router)
	sellerHandler.Configure(router)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router, nil
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic occurred:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
