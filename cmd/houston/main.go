package main

import (
	"github.com/ternaryss/houston/internal/app/db"
	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/web"
)

func main() {
	settings := settings.LoadSettings()
	dbProvider := db.NewDbProvider(settings)
	defer dbProvider.CloseConnection()
	usersStore := db.NewUsersStore(dbProvider.Db())
	errorsHandler := handlers.NewErrorsHandler()
	dashboardHandler := handlers.NewDashboardHandler()
	usersHandler := handlers.NewUsersHandler(usersStore)
	server := web.NewServer(settings, errorsHandler, dashboardHandler, usersHandler)
	server.Run()
}
