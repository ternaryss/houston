package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type deleteWebAppCmd struct {
	webAppsStore      types.WebAppsStore
	subscribersStore  types.SubscribersStore
	healthChecksStore types.HealthChecksStore
}

func NewDeleteWebAppCmd(was types.WebAppsStore, sus types.SubscribersStore, hcs types.HealthChecksStore) *deleteWebAppCmd {
	return &deleteWebAppCmd{
		webAppsStore:      was,
		subscribersStore:  sus,
		healthChecksStore: hcs,
	}
}

func (c *deleteWebAppCmd) Execute(id, usr string) error {
	slog.Info("Deleting web application", "id", id, "user", usr)

	if id == "" {
		return sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(id, usr, nil)

	if err != nil {
		return err
	}

	tx, err := c.webAppsStore.Begin()

	if err != nil {
		return err
	}

	if err := c.healthChecksStore.DeleteByWebAppId(app.Id, tx); err != nil {
		c.webAppsStore.Rollback(tx)
		return err
	}

	if err := c.subscribersStore.DeleteByWebAppId(app.Id, tx); err != nil {
		c.webAppsStore.Rollback(tx)
		return err
	}

	if err := c.webAppsStore.DeleteByIdAndUserEmail(app.Id, app.UserEmail, tx); err != nil {
		c.webAppsStore.Rollback(tx)
		return err
	}

	if err := c.webAppsStore.Commit(tx); err != nil {
		return err
	}

	slog.Info("Web application deleted", "id", app.Id)

	return nil
}
