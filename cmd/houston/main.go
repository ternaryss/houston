package main

import (
	"github.com/ternaryss/houston/internal/app/db"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/web"
)

func main() {
	settings := settings.LoadSettings()
	dbProvider := db.NewDbProvider(settings)
	defer dbProvider.CloseConnection()
	server := web.NewServer(settings)
	server.Run()
}
