package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/settings"
)

type server struct {
	settings         settings.Settings
	errorsHandler    *handlers.ErrorsHandler
	dashboardHandler *handlers.DashboardHandler
	usersHandler     *handlers.UsersHandler
}

func NewServer(
	stg settings.Settings,
	erh *handlers.ErrorsHandler,
	dsh *handlers.DashboardHandler,
	ush *handlers.UsersHandler,
) *server {
	return &server{
		settings:         stg,
		errorsHandler:    erh,
		dashboardHandler: dsh,
		usersHandler:     ush,
	}
}

func (s *server) middlewares() middleware {
	return middlewaresChain(
		requestTimeMiddleware,
	)
}

func (s *server) staticFiles(rtr *http.ServeMux) {
	public := http.FileServer(http.Dir("public"))
	rtr.Handle("/public/", http.StripPrefix("/public/", public))
}

func (s *server) routes(rtr *http.ServeMux) {
	rtr.HandleFunc("/not-found", s.errorsHandler.NotFoundError)
	rtr.HandleFunc("/error", s.errorsHandler.InternalServerError)
	rtr.HandleFunc("/sign-up", s.usersHandler.SignUp)
	rtr.HandleFunc("/{$}", s.dashboardHandler.Dashboard)
	rtr.HandleFunc("/", s.errorsHandler.NotFoundError)
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", s.settings.Server.Host, s.settings.Server.Port)
	router := http.NewServeMux()
	s.staticFiles(router)
	s.routes(router)
	server := &http.Server{
		Addr:    addr,
		Handler: s.middlewares()(router),
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
