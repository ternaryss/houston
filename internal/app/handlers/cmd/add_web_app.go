package cmd

import (
	"log/slog"
	"strings"

	"github.com/ternaryss/houston/internal/app/types"
)

type addWebAppCmd struct {
	webAppsStore     types.WebAppsStore
	subscribersStore types.SubscribersStore
}

func NewAddWebAppCmd(was types.WebAppsStore, sus types.SubscribersStore) *addWebAppCmd {
	return &addWebAppCmd{
		webAppsStore:     was,
		subscribersStore: sus,
	}
}

func (c *addWebAppCmd) Execute(frm types.WebAppForm, usr string) (string, error) {
	slog.Info("Adding web application", "form", frm, "user", usr)
	frm.Validate()

	if len(frm.Errors) > 0 {
		slog.Info("Add web application form not valid", "errors", frm.Errors)
		return "", nil
	}

	tx, err := c.webAppsStore.Begin()

	if err != nil {
		return "", nil
	}

	app := types.NewWebApp(frm.Name, frm.Url, frm.Interval, usr, frm.Status)
	app, err = c.webAppsStore.Insert(app, tx)

	if err != nil {
		c.webAppsStore.Rollback(tx)
		return "", err
	}

	for _, email := range frm.Notify {
		subscriber := types.NewSubscriber(app.Id, strings.ToLower(email))

		if _, err := c.subscribersStore.Insert(subscriber, tx); err != nil {
			c.webAppsStore.Rollback(tx)
			return "", err
		}
	}

	if err := c.webAppsStore.Commit(tx); err != nil {
		return "", err
	}

	slog.Info("New web app added", "app", app)

	return app.Id, nil
}
