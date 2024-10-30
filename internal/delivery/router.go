package delivery

import (
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
)

func NewRouter(cfg config.Config) *mux.Router {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		return nil
	}
	rdb, err := connector.GetRedisConnector(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil
	}

	userRepo := postgres.NewUserRepository(dbPool)
	userUC := service.NewUserService(userRepo, zap.L())

	sessionRepo := redis.NewSessionRepository(rdb, int(cfg.Session.ExpirationTime.Seconds()), zap.L())
	sessionUC := service.NewAuthService(sessionRepo, zap.L())

	sessionManager := utils.NewSessionManager(sessionUC, int(cfg.Session.ExpirationTime.Seconds()), cfg.Session.SecureCookie, zap.L())

	authHandler := http3.NewAuthEndpoints(sessionUC, sessionManager, zap.L())
	userHandler := http3.NewUserEndpoints(userUC, sessionUC, sessionManager, zap.L())

	authHandler.Configure(router)
	userHandler.Configure(router)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
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

// func isAuthenticated(r *http.Request, authHandler *http3.AuthHandler) bool {
// 	cookie, err := r.Cookie("session_id")
// 	if err != nil || cookie == nil {
// 		log.Println("No session cookie found")

// 		return false
// 	}

// 	exists := authHandler.SessionRepo.SessionExists(cookie.Value)
// 	log.Printf("Session exists: %v for session_id: %s", exists, cookie.Value)

// 	return exists
// }

// func authMiddleware(authHandler *http3.AuthHandler) mux.MiddlewareFunc {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if isAuthenticated(r, authHandler) {
// 				w.Header().Set("X-Authenticated", "true")
// 			} else {
// 				w.Header().Set("X-Authenticated", "false")
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
