package cmd

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/ternaryss/houston/internal/app/types"
)

type editWebAppCmd struct {
	webAppsStore     types.WebAppsStore
	subscribersStore types.SubscribersStore
}

func NewEditWebAppCmd(was types.WebAppsStore, sus types.SubscribersStore) *editWebAppCmd {
	return &editWebAppCmd{
		webAppsStore:     was,
		subscribersStore: sus,
	}
}

func (c *editWebAppCmd) Execute(frm types.WebAppForm, id, usr string) error {
	slog.Info("Editing web application", "id", id, "form", frm, "user", usr)

	if id == "" {
		return sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(id, usr)

	if err != nil {
		return err
	}

	frm.Validate()

	if len(frm.Errors) > 0 {
		slog.Info("Edit web application form not valid", "errors", frm.Errors)
		return nil
	}

	app.Name = frm.Name
	app.Url = frm.Url
	app.UserEmail = usr

	tx, err := c.webAppsStore.Begin()

	if err != nil {
		return err
	}

	if err := c.subscribersStore.DeleteByWebAppId(app.Id); err != nil {
		c.webAppsStore.Rollback(tx)
		return err
	}

	if _, err := c.webAppsStore.Update(app); err != nil {
		c.webAppsStore.Rollback(tx)
		return err
	}

	for _, email := range frm.Notify {
		subscriber := types.NewSubscriber(app.Id, strings.ToLower(email))

		if _, err := c.subscribersStore.Insert(subscriber); err != nil {
			c.webAppsStore.Rollback(tx)
			return err
		}
	}

	if err := c.webAppsStore.Commit(tx); err != nil {
		return err
	}

	slog.Info("Web application modified", "app", app)

	return nil
}
