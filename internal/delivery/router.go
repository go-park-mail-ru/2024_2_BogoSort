package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	http3 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

func handleRepoError(err error, message string) error {
	zap.L().Error(message, zap.Error(err))
	return errors.Wrap(err, message)
}

func NewRouter(cfg config.Config) (*mux.Router, error) {
	log := zap.Must(zap.NewProduction())
	defer log.Sync()

	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	policy := bluemonday.UGCPolicy()

	authRouter := router.PathPrefix("").Subrouter()

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		log.Error("Failed to connect to Postgres", zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to Postgres")
	}
	rdb, err := connector.GetRedisConnector(cfg.RdAddr, cfg.RdPass, cfg.RdDB)
	if err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	ctx := context.Background()

	advertsRepo, err := postgres.NewAdvertRepository(dbPool, log, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create advert repository")
	}
	staticRepo, err := postgres.NewStaticRepository(ctx, dbPool, cfg.Static.Path, cfg.Static.MaxSize, log, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create static repository")
	}
	categoryRepo, err := postgres.NewCategoryRepository(dbPool, log, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create category repository")
	}
	sessionRepo, err := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), ctx, log)
	if err != nil {
		return nil, handleRepoError(err, "unable to create session repository")
	}
	userRepo, err := postgres.NewUserRepository(dbPool, ctx, log, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create user repository")
	}
	sellerRepo, err := postgres.NewSellerRepository(dbPool, ctx, log)
	if err != nil {
		return nil, handleRepoError(err, "unable to create seller repository")
	}
	cartRepo, err := postgres.NewCartRepository(dbPool, ctx, log)
	if err != nil {
		return nil, handleRepoError(err, "unable to create cart repository")
	}
	purchaseRepo, err := postgres.NewPurchaseRepository(dbPool, log, ctx, cfg.PGTimeout)
	if err != nil {
		return nil, handleRepoError(err, "unable to create purchase repository")
	}
	csrfToken, err := utils.NewAesCryptHashToken(log)
	if err != nil {
		return nil, handleRepoError(err, "unable to create csrf token")
	}

	advertsUC := service.NewAdvertService(advertsRepo, staticRepo, sellerRepo, log)
	staticUC := service.NewStaticService(staticRepo, log)
	categoryUC := service.NewCategoryService(categoryRepo, log)
	purchaseUC := service.NewPurchaseService(purchaseRepo, advertsRepo, cartRepo, log)
	cartUC := service.NewCartService(cartRepo, advertsRepo, log)
	userUC := service.NewUserService(userRepo, sellerRepo, log)
	sessionUC := service.NewAuthService(sessionRepo, log)
	sessionManager := utils.NewSessionManager(sessionUC, int(cfg.Session.ExpirationTime.Seconds()), cfg.Session.SecureCookie, log)

	router.Use(middleware.NewAuthMiddleware(sessionManager).AuthMiddleware)

	advertsHandler := http3.NewAdvertEndpoints(advertsUC, staticUC, sessionManager, log, policy)
	authHandler := http3.NewAuthEndpoints(sessionUC, sessionManager, log)
	userHandler := http3.NewUserEndpoints(userUC, sessionUC, sessionManager, staticUC, log, policy)
	sellerHandler := http3.NewSellerEndpoints(sellerRepo, log)
	purchaseHandler := http3.NewPurchaseEndpoints(purchaseUC, log)
	cartHandler := http3.NewCartEndpoints(cartUC, log)
	categoryHandler := http3.NewCategoryEndpoints(categoryUC, log)
	staticHandler := http3.NewStaticEndpoints(staticUC, log)

	csrfEndpoints := http3.NewCSRFEndpoints(csrfToken, sessionManager)
	csrfEndpoints.Configure(router)
	userHandler.ConfigureUnprotectedRoutes(router)

	authRouter.Use(middleware.CSRFMiddleware(csrfToken, sessionManager))

	advertsHandler.ConfigureRoutes(authRouter)
	categoryHandler.ConfigureRoutes(authRouter)
	authHandler.Configure(authRouter)
	userHandler.ConfigureProtectedRoutes(authRouter)
	sellerHandler.Configure(authRouter)
	cartHandler.Configure(authRouter)
	staticHandler.ConfigureRoutes(authRouter)
	purchaseHandler.ConfigureRoutes(authRouter)
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
