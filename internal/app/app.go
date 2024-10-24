package app

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	if err := config.ServerInit(); err != nil {
		zap.L().Error("config init error", zap.Error(err))
		return err
	}

	router := delivery.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://two024-2-bogo-sort.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	zap.L().Info("server started", zap.String("address", config.GetServerAddress()))

	server.server = &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      corsHandler,
		ReadTimeout:  config.GetReadTimeout(),
		WriteTimeout: config.GetWriteTimeout(),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	err := server.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		zap.L().Error("server failed", zap.Error(err))
		return errors.Wrap(err, "server failed")
	}

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
