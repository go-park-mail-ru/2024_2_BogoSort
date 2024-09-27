package app

import (
	"net/http"
	"time"
	"context"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run() error {
	srv.server = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
