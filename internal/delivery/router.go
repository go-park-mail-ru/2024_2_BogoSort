package delivery

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	service "github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/services"
	http3 "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func NewRouter(cfg config.Config) *mux.Router {
    router := mux.NewRouter()
    router.Use(recoveryMiddleware)

    dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
    if err != nil {
        zap.L().Error("unable to connect to database", zap.Error(err))
        return nil
    }

    advertsRepo, err := postgres.NewAdvertRepository(dbPool, zap.L())
    if err != nil {
        zap.L().Error("unable to create advert repository", zap.Error(err))
        return nil
    }

    advertsUseCase, err := service.NewAdvertService(advertsRepo, zap.L())
    if err != nil {
        zap.L().Error("unable to create advert use case", zap.Error(err))
        return nil
    }

    advertsHandler := http3.NewAdvertEndpoints(advertsUseCase, zap.L())

    advertsHandler.ConfigureRoutes(router)

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

/*func isAuthenticated(r *http.Request, authHandler *http3.AuthHandler) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie == nil {
		log.Println("No session cookie found")

		return false
	}

	exists := authHandler.SessionRepo.SessionExists(cookie.Value)
	log.Printf("Session exists: %v for session_id: %s", exists, cookie.Value)

	return exists
}

func authMiddleware(authHandler *http3.AuthHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAuthenticated(r, authHandler) {
				w.Header().Set("X-Authenticated", "true")
			} else {
				w.Header().Set("X-Authenticated", "false")
			}

			next.ServeHTTP(w, r)
		})
	}
}*/
