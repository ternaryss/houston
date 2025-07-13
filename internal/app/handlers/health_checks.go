package handlers

import (
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
)

type HealthChecksHandler struct {
	settings          settings.Settings
	webAppsStore      types.WebAppsStore
	healthChecksStore types.HealthChecksStore
}

func NewHealthChecksHandler(stg settings.Settings, was types.WebAppsStore, hcs types.HealthChecksStore) *HealthChecksHandler {
	return &HealthChecksHandler{
		settings:          stg,
		webAppsStore:      was,
		healthChecksStore: hcs,
	}
}

func (h *HealthChecksHandler) PeriodicHealthCheck(itv string) {
	cmd.NewPeriodicHealthCheckCmd(h.webAppsStore, h.healthChecksStore).Execute(itv)
}

func (h *HealthChecksHandler) HealthChecksRetention() {
	cmd.NewRetentionCmd(h.healthChecksStore).Execute(h.settings.Retention.OlderThan)
}
