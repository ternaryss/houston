package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type editWebAppCmd struct {
	webAppsStore types.WebAppsStore
}

func NewEditWebAppCmd(was types.WebAppsStore) *editWebAppCmd {
	return &editWebAppCmd{
		webAppsStore: was,
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

	if _, err := c.webAppsStore.Update(app); err != nil {
		return err
	}

	slog.Info("Web application modified", "app", app)

	return nil
}
