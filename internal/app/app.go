package app

import (
	"net/http"
	"time"

	"emporium/internal/myhandlers"
)

const (
	Timeout = time.Second * 3
	Address = ":8080"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run() error {
	http.HandleFunc("/adverts", myhandlers.GetadvertsHandler)

	return srv.server.ListenAndServe()
}
