package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type periodicHealthCheckCmd struct {
	timeout           int
	concurrent        int
	webAppsStore      types.WebAppsStore
	subscribersStore  types.SubscribersStore
	healthChecksStore types.HealthChecksStore
	emailClient       types.EmailClient
}

func NewPeriodicHealthCheckCmd(
	was types.WebAppsStore,
	sus types.SubscribersStore,
	hcs types.HealthChecksStore,
	emc types.EmailClient,
) *periodicHealthCheckCmd {
	return &periodicHealthCheckCmd{
		timeout:           5,
		concurrent:        10,
		webAppsStore:      was,
		subscribersStore:  sus,
		healthChecksStore: hcs,
		emailClient:       emc,
	}
}

func (c *periodicHealthCheckCmd) notify(app *types.WebApp, acl, lst *types.HealthCheck) {
	var notification types.Message
	slog.Info("Notifying about web application health", "app", app, "actual", acl, "last", lst)

	if lst != nil && lst.Status == acl.Status {
		slog.Info("Notifying aborted, no status change", "appId", app.Id)
		return
	}

	if acl.Status == app.Status {
		notification = types.NewAliveMessage(app, acl)
	} else {
		notification = types.NewNotAliveMessage(app, acl)
	}

	subscribers, err := c.subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id, nil)

	if err != nil {
		slog.Error("Fetching subscribers failed", "appId", app.Id, "err", err)
	}

	for _, subscriber := range subscribers {
		if err := c.emailClient.SendEmail(subscriber.Email, notification.Subject, notification.Content); err != nil {
			slog.Error("Sending notification failed", "appId", app.Id, "subscriber", subscriber.Email, "err", err)
		}
	}
}

func (c *periodicHealthCheckCmd) checkHealth(ctx context.Context, app *types.WebApp) {
	slog.Info("Executing health check", "app", app)
	httpStatus := -1
	status := types.HealthyOk
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(c.timeout)*time.Second)
	defer cancel()
	latest, err := c.healthChecksStore.GetFirstByWebAppIdOrderByCreatedAtDesc(app.Id, nil)

	if err != nil && err != sql.ErrNoRows {
		slog.Error("Fetching latest health check failed", "appId", app.Id, "err", err)
		return
	}

	if request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, app.Url, nil); err != nil {
		slog.Error("Creation of HTTP request failed", "appId", app.Id, "err", err)
		status = types.HealthyErr
	} else {
		if response, err := http.DefaultClient.Do(request); err != nil {
			slog.Error("GET request to web application failed", "appId", app.Id, "err", err)
			status = types.HealthyErr
		} else {
			defer response.Body.Close()

			if response.StatusCode != app.Status {
				status = types.HealthyErr
			}

			httpStatus = response.StatusCode
		}
	}

	app.Healthy = status
	health := types.NewHealthCheck(app.Id, httpStatus)
	tx, err := c.healthChecksStore.Begin()

	if err != nil {
		slog.Error("Begin of transaction failed", "appId", app.Id, "err", err)
		return
	}

	if _, err := c.healthChecksStore.Insert(health, tx); err != nil {
		slog.Error("Health check insert failed", "appId", app.Id, "err", err)
		c.healthChecksStore.Rollback(tx)
		return
	}

	if _, err := c.webAppsStore.Update(app, tx); err != nil {
		slog.Error("Web application health status update failed", "appId", app.Id, "err", err)
		c.healthChecksStore.Rollback(tx)
		return
	}

	if err := c.healthChecksStore.Commit(tx); err != nil {
		slog.Error("Commit of transaction failed", "appId", app.Id, "err", err)
	}

	c.notify(app, health, latest)
	slog.Info("Health check executed", "app", app)
}

func (c *periodicHealthCheckCmd) Execute(itv string) {
	slog.Info("Checking health", "interval", itv)
	apps, err := c.webAppsStore.GetByIntervalOrderByNameAsc(itv, nil)

	if err != nil {
		slog.Error("Fetching web applications failed", "err", err)
		return
	}

	slog.Info("Web applications fetched for health check", "interval", itv, "apps", len(apps))
	queue := make(chan struct{}, c.concurrent)
	done := make(chan struct{})

	for _, app := range apps {
		queue <- struct{}{}

		go func(app *types.WebApp) {
			defer func() { <-queue; done <- struct{}{} }()
			c.checkHealth(context.Background(), app)
		}(app)
	}

	for range apps {
		<-done
	}

	slog.Info("Health checking completed", "interval", itv)
}
