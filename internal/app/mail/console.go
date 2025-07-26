package mail

import (
	"log/slog"

	"github.com/ternaryss/houston/internal/app/settings"
)

type consoleClient struct {
	settings settings.Settings
}

func NewConsoleClient(stg settings.Settings) *consoleClient {
	return &consoleClient{
		settings: stg,
	}
}

func (c *consoleClient) SendEmail(to, sub, msg string) error {
	conf := c.settings.Smtp
	from := conf.User

	if conf.From != "" {
		from = conf.From
	}

	slog.Info("E-mail sent", "user", conf.User, "from", from, "to", to, "subject", sub, "message", msg)

	return nil
}
