package cron

import (
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/types"
)

type fiveMinutesIntervalScheduler struct {
	cron                string
	interval            string
	healthChecksHandler *handlers.HealthChecksHandler
}

func NewFiveMinutesIntervalScheduler(hch *handlers.HealthChecksHandler) *fiveMinutesIntervalScheduler {
	return &fiveMinutesIntervalScheduler{
		cron:                "*/5 * * * *",
		interval:            types.Interval5M,
		healthChecksHandler: hch,
	}
}

func (s *fiveMinutesIntervalScheduler) Run() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		slog.Error("Five minutes interval scheduler failed", "err", err)
		os.Exit(1)
	}

	job := gocron.CronJob(s.cron, false)
	task := gocron.NewTask(func() {
		s.healthChecksHandler.PeriodicHealthCheck(s.interval)
	})

	if _, err := scheduler.NewJob(job, task); err != nil {
		slog.Error("Five minutes interval job failed", "err", err)
		os.Exit(1)
	}

	scheduler.Start()
	slog.Info("Five minutes interval scheduler configured", "interval", s.interval, "cron", s.cron)
}
