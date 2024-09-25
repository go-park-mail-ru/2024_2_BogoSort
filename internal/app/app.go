package app

import (
	"net/http"
	"time"

	"emporium/internal/handlers"
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
	handlers.AddTestAdvert()
	http.HandleFunc("/adverts", handlers.GetAdvertsHandler)

	return srv.server.ListenAndServe()
}
