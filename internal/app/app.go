package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/handlers"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/utils"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	if err := config.Init(); err != nil {
		return err
	}
	utils.InitJWT()

	router := handlers.NewRouter()

	s.server = &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}