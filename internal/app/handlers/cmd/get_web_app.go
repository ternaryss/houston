package cmd

import (
	"database/sql"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type getWebAppCmd struct {
	webAppsStore     types.WebAppsStore
	subscribersStore types.SubscribersStore
}

func NewGetWebAppCmd(was types.WebAppsStore, sus types.SubscribersStore) *getWebAppCmd {
	return &getWebAppCmd{
		webAppsStore:     was,
		subscribersStore: sus,
	}
}

func (c *getWebAppCmd) Execute(id, usr string) (*types.WebApp, []*types.Subscriber, error) {
	slog.Info("Fetching web application", "id", id, "user", usr)

	if id == "" {
		return nil, []*types.Subscriber{}, sql.ErrNoRows
	}

	app, err := c.webAppsStore.GetByIdAndUserEmail(id, usr)

	if err != nil {
		return nil, []*types.Subscriber{}, err
	}

	subscribers, err := c.subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id)

	if err != nil {
		return nil, []*types.Subscriber{}, err
	}

	slog.Info("Web application fetched", "app", app, "subscribers", len(subscribers))

	return app, subscribers, nil
}
