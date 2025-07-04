package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type deleteWebAppCmd struct {
	webAppsStore types.WebAppsStore
}

func NewDeleteWebAppCmd(was types.WebAppsStore) *deleteWebAppCmd {
	return &deleteWebAppCmd{
		webAppsStore: was,
	}
}

func (c *deleteWebAppCmd) Execute(id, usr string) error {
	slog.Info("Deleting web application", "id", id, "user", usr)

	if id == "" {
		return sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(id, usr)

	if err != nil {
		return err
	}

	if err := c.webAppsStore.DeleteByIdAndUserEmail(app.Id, app.UserEmail); err != nil {
		return err
	}

	slog.Info("Web application deleted", "id", app.Id)

	return nil
}
