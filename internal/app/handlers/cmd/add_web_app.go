package cmd

import (
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type addWebAppCmd struct {
	webAppsStore types.WebAppsStore
}

func NewAddWebAppCmd(was types.WebAppsStore) *addWebAppCmd {
	return &addWebAppCmd{
		webAppsStore: was,
	}
}

func (c *addWebAppCmd) Execute(frm types.WebAppForm, usr string) (string, error) {
	slog.Info("Adding web application", "form", frm, "user", usr)
	frm.Validate()

	if len(frm.Errors) > 0 {
		slog.Info("Add web application form not valid", "errors", frm.Errors)
		return "", nil
	}

	app := types.NewWebApp(frm.Name, frm.Url, usr)

	if _, err := c.webAppsStore.Insert(app); err != nil {
		return "", err
	}

	slog.Info("New web app added", "app", app)

	return app.Id, nil
}
