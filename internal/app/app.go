package app

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/handlers"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run() error {
	router := handlers.NewRouter()
	srv.server = &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.server.ListenAndServe()
}
