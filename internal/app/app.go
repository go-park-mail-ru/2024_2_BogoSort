package app

import (
	"context"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/handlers"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Run() error {
	if err := config.Init(); err != nil {
		return err
	}

	utils.InitJWT()

	router := handlers.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://two024-2-bogo-sort.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	log.Printf("Server started on %s", config.GetServerAddress())

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
		log.Fatalf("server failed: %v", err)
	}

	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownTimeout())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("server exiting")

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	err := server.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
