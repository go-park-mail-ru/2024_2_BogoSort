package app

import (
	"net/http"
	"time"

	"emporium/internal/handlers"
)

type Server struct {
	server *http.Server
}

func (s *Server) AdvertsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetAdvertsHandler(w, r)
	case http.MethodPost:
		handlers.AddAdvertHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (srv *Server) Run() error {
	srv.server = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc("/adverts", srv.AdvertsHandler)
	return srv.server.ListenAndServe()
}
