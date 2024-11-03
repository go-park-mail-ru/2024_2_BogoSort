package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
}

func (server *Server) Run() error {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	cfg, err := config.Init()
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	router, err := delivery.NewRouter(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to initialize router")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://two024-2-bogo-sort.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}).Handler(router)

	zap.L().Info("Server started on " + config.GetServerAddress())

	server.server = &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      corsHandler,
		ReadTimeout:  config.GetReadTimeout(),
		WriteTimeout: config.GetWriteTimeout(),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-stop
	zap.L().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownTimeout())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("server exiting")

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	err := server.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
