package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/settings"
)

type server struct {
	settings            settings.Settings
	errorsHandler       *handlers.ErrorsHandler
	dashboardHandler    *handlers.DashboardHandler
	webAppsHandler      *handlers.WebAppsHandler
	healthChecksHandler *handlers.HealthChecksHandler
	usersHandler        *handlers.UsersHandler
}

func NewServer(
	stg settings.Settings,
	erh *handlers.ErrorsHandler,
	dsh *handlers.DashboardHandler,
	wah *handlers.WebAppsHandler,
	hch *handlers.HealthChecksHandler,
	ush *handlers.UsersHandler,
) *server {
	return &server{
		settings:            stg,
		errorsHandler:       erh,
		dashboardHandler:    dsh,
		webAppsHandler:      wah,
		healthChecksHandler: hch,
		usersHandler:        ush,
	}
}

func (s *server) middlewares() middleware {
	return middlewaresChain(
		authorizationMiddleware(s.settings),
		requestTimeMiddleware,
	)
}

func (s *server) staticFiles(rtr *http.ServeMux) {
	public := http.FileServer(http.Dir("public"))
	rtr.Handle("/public/", http.StripPrefix("/public/", public))
}

func (s *server) routes(rtr *http.ServeMux, ebd bool) {
	rtr.HandleFunc("/not-found", s.errorsHandler.NotFoundError)
	rtr.HandleFunc("/error", s.errorsHandler.InternalServerError)
	rtr.HandleFunc("GET /sign-in", s.usersHandler.SignIn)
	rtr.HandleFunc("POST /sign-in", s.usersHandler.SignIn)
	rtr.HandleFunc("POST /sign-out", s.usersHandler.SignOut)

	if ebd {
		rtr.HandleFunc("GET /sign-up", s.usersHandler.SignUp)
		rtr.HandleFunc("POST /sign-up", s.usersHandler.SignUp)
	}

	rtr.HandleFunc("GET /web-apps", s.webAppsHandler.GetWebApps)
	rtr.HandleFunc("GET /web-apps/add", s.webAppsHandler.AddWebApp)
	rtr.HandleFunc("POST /web-apps", s.webAppsHandler.AddWebApp)
	rtr.HandleFunc("GET /web-apps/{id}", s.webAppsHandler.GetWebApp)
	rtr.HandleFunc("GET /web-apps/{id}/edit", s.webAppsHandler.EditWebApp)
	rtr.HandleFunc("PUT /web-apps/{id}", s.webAppsHandler.EditWebApp)
	rtr.HandleFunc("DELETE /web-apps/{id}", s.webAppsHandler.DeleteWebApp)
	rtr.HandleFunc("GET /web-apps/{id}/health", s.healthChecksHandler.GetHealthChecks)
	rtr.HandleFunc("GET /{$}", s.dashboardHandler.Dashboard)
	rtr.HandleFunc("/", s.errorsHandler.NotFoundError)
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", s.settings.Server.Host, s.settings.Server.Port)
	router := http.NewServeMux()
	s.staticFiles(router)
	s.routes(router, s.settings.Authorization.SignUp.Enabled)
	server := &http.Server{
		Addr:    addr,
		Handler: s.middlewares()(router),
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
