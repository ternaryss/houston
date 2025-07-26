package cmd

import (
	"log/slog"

	"github.com/ternaryss/houston/internal/app/types"
)

type getDoughnutHealthChartCmd struct {
	healthChecksStore types.HealthChecksStore
}

func NewGetDoughnutHealthChartCmd(hcs types.HealthChecksStore) *getDoughnutHealthChartCmd {
	return &getDoughnutHealthChartCmd{
		healthChecksStore: hcs,
	}
}

func (c *getDoughnutHealthChartCmd) Execute(app *types.WebApp) (types.HealthStats, error) {
	slog.Info("Fetching doughnut health chart", "webAppId", app.Id)
	up, err := c.healthChecksStore.CountByWebAppIdAndStatus(app.Id, app.Status, nil)

	if err != nil {
		return types.HealthStats{}, err
	}

	down, err := c.healthChecksStore.CountByWebAppIdAndNotStatus(app.Id, app.Status, nil)

	if err != nil {
		return types.HealthStats{}, err
	}

	slog.Info("Doughnut health chart fetched", "webAppId", app.Id, "up", up, "down", down)

	return types.NewHealthStats(up, down), nil
}
