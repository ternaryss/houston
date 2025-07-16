package cron

import (
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/ternaryss/houston/internal/app/handlers"
	"github.com/ternaryss/houston/internal/app/types"
)

type fifteenMinutesIntervalScheduler struct {
	cron                string
	interval            string
	healthChecksHandler *handlers.HealthChecksHandler
}

func NewFifteenMinutesIntervalScheduler(hch *handlers.HealthChecksHandler) *fifteenMinutesIntervalScheduler {
	return &fifteenMinutesIntervalScheduler{
		cron:                "*/15 * * * *",
		interval:            types.Interval15M,
		healthChecksHandler: hch,
	}
}

func (s *fifteenMinutesIntervalScheduler) Run() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		slog.Error("Fifteen minutes interval scheduler failed", "err", err)
		os.Exit(1)
	}

	job := gocron.CronJob(s.cron, false)
	task := gocron.NewTask(func() {
		s.healthChecksHandler.PeriodicHealthCheck(s.interval)
	})

	if _, err := scheduler.NewJob(job, task); err != nil {
		slog.Error("Fifteen minutes interval job failed", "err", err)
		os.Exit(1)
	}

	scheduler.Start()
	slog.Info("Fifteen minutes interval scheduler configured", "interval", s.interval, "cron", s.cron)
}
