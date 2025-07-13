package cmd

import (
	"log/slog"
	"time"

	"github.com/ternaryss/houston/internal/app/types"
)

type retentionCmd struct {
	healthChecksStore types.HealthChecksStore
}

func NewRetentionCmd(hcs types.HealthChecksStore) *retentionCmd {
	return &retentionCmd{
		healthChecksStore: hcs,
	}
}

func (c *retentionCmd) Execute(day int) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	days := time.Hour * 24 * time.Duration(day)
	olderThan := today.Add(-days)
	slog.Info("Executing health checks retention", "olderThan", olderThan)

	if err := c.healthChecksStore.DeleteByCreatedAtLowerThan(olderThan); err != nil {
		slog.Error("Health checks retention failed", "err", err)
		return
	}

	slog.Info("Health checks retention executed", "olderThan", olderThan)
}
