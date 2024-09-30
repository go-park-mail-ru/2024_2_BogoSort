package app

import (
	"context"
	"net/http"
	"time"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/handlers"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
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

	server.server = &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      corsHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.server.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "failed to listen and serve")
	}
	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	err := server.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}
	return nil
}
