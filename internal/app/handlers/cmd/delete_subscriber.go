package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type deleteSubscriberCmd struct {
	webAppsStore     types.WebAppsStore
	subscribersStore types.SubscribersStore
}

func NewDeleteSubscriberCmd(was types.WebAppsStore, sus types.SubscribersStore) *deleteSubscriberCmd {
	return &deleteSubscriberCmd{
		webAppsStore:     was,
		subscribersStore: sus,
	}
}

func (c *deleteSubscriberCmd) Execute(wid, usr string) error {
	slog.Info("Deleting web application subscriber", "webAppId", wid, "subscriber", usr)

	if wid == "" {
		return sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(wid, usr, nil)

	if err != nil {
		return err
	}

	if err := c.subscribersStore.DeleteByWebAppIdAndEmail(app.Id, usr, nil); err != nil {
		return err
	}

	slog.Info("Web application subscriber deleted", "webAppId", app.Id, "subscriber", usr)

	return nil
}
