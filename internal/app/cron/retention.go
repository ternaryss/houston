package cron

import (
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/settings"
)

type retentionScheduler struct {
	cron                string
	settings            settings.Settings
	healthChecksHandler *handlers.HealthChecksHandler
}

func NewRetentionScheduler(stg settings.Settings, hch *handlers.HealthChecksHandler) *retentionScheduler {
	return &retentionScheduler{
		cron:                "0 0 * * *",
		settings:            stg,
		healthChecksHandler: hch,
	}
}

func (s *retentionScheduler) Run() {
	enabled := s.settings.Retention.Enabled

	if enabled {
		scheduler, err := gocron.NewScheduler()

		if err != nil {
			slog.Error("Retention scheduler failed", "err", err)
			os.Exit(1)
		}

		job := gocron.CronJob(s.cron, false)
		task := gocron.NewTask(func() {
			s.healthChecksHandler.HealthChecksRetention()
		})

		if _, err := scheduler.NewJob(job, task); err != nil {
			slog.Error("Retention job failed", "err", err)
			os.Exit(1)
		}

		scheduler.Start()
	}

	slog.Info("Retention scheduler configured", "enabled", enabled, "cron", s.cron)
}
