package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type getWebAppCmd struct {
	webAppsStore types.WebAppsStore
}

func NewGetWebAppCmd(was types.WebAppsStore) *getWebAppCmd {
	return &getWebAppCmd{
		webAppsStore: was,
	}
}

func (c *getWebAppCmd) Execute(id, usr string) (*types.WebApp, error) {
	slog.Info("Fetching web application", "id", id, "user", usr)

	if id == "" {
		return nil, sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(id, usr)

	if err != nil {
		return nil, err
	}

	slog.Info("Web application fetched", "app", app)

	return app, nil
}
