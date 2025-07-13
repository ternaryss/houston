package cron

import (
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/types"
)

type oneHourIntervalScheduler struct {
	cron                string
	interval            string
	healthChecksHandler *handlers.HealthChecksHandler
}

func NewOneHourIntervalScheduler(hch *handlers.HealthChecksHandler) *oneHourIntervalScheduler {
	return &oneHourIntervalScheduler{
		cron:                "0 * * * *",
		interval:            types.Interval1H,
		healthChecksHandler: hch,
	}
}

func (s *oneHourIntervalScheduler) Run() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		slog.Error("One hour interval scheduler failed", "err", err)
		os.Exit(1)
	}

	job := gocron.CronJob(s.cron, false)
	task := gocron.NewTask(func() {
		s.healthChecksHandler.PeriodicHealthCheck(s.interval)
	})

	if _, err := scheduler.NewJob(job, task); err != nil {
		slog.Error("One hour interval job failed", "err", err)
		os.Exit(1)
	}

	scheduler.Start()
	slog.Info("One hour interval scheduler configured", "interval", s.interval, "cron", s.cron)
}
