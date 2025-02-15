package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/settings"
)

type server struct {
	settings settings.Settings
}

func NewServer(stg settings.Settings) *server {
	return &server{
		settings: stg,
	}
}

func (s *server) staticFiles(rtr *http.ServeMux) {
	public := http.FileServer(http.Dir("public"))
	rtr.Handle("/public/", http.StripPrefix("/public/", public))
}

func (s *server) routes(rtr *http.ServeMux) {
	rtr.HandleFunc("/", handlers.Hello)
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", s.settings.Server.Host, s.settings.Server.Port)
	router := http.NewServeMux()
	s.staticFiles(router)
	s.routes(router)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
