package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views/components"
)

type HealthChecksHandler struct {
	settings          settings.Settings
	webAppsStore      types.WebAppsStore
	subscribersStore  types.SubscribersStore
	healthChecksStore types.HealthChecksStore
}

func NewHealthChecksHandler(
	stg settings.Settings,
	was types.WebAppsStore,
	sus types.SubscribersStore,
	hcs types.HealthChecksStore,
) *HealthChecksHandler {
	return &HealthChecksHandler{
		settings:          stg,
		webAppsStore:      was,
		subscribersStore:  sus,
		healthChecksStore: hcs,
	}
}

func (h *HealthChecksHandler) GetHealthChecks(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	id := req.PathValue("id")
	app, _, err := cmd.NewGetWebAppCmd(h.webAppsStore, h.subscribersStore).Execute(id, user)

	if err != nil {
		if err == sql.ErrNoRows {
			helpers.NotFoundError(res, req)
			return
		}

		helpers.InternalServerError(err, res, req)
		return
	}

	page, err := cmd.NewGetHealthChecksCmd(h.healthChecksStore).Execute(id, req.URL.Query())

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	template := components.HealthChecksList(app, page)
	helpers.Render(template, res, req)
}

func (h *HealthChecksHandler) PeriodicHealthCheck(itv string) {
	cmd.NewPeriodicHealthCheckCmd(h.webAppsStore, h.healthChecksStore).Execute(itv)
}

func (h *HealthChecksHandler) HealthChecksRetention() {
	cmd.NewRetentionCmd(h.healthChecksStore).Execute(h.settings.Retention.OlderThan)
}
