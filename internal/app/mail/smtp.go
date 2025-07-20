package mail

import (
	"crypto/tls"
	"log/slog"

	"github.com/ternaryss/houston/internal/app/settings"
	"gopkg.in/gomail.v2"
)

type smtpClient struct {
	settings settings.Settings
	dialer   *gomail.Dialer
}

func NewSmtpClient(stg settings.Settings) *smtpClient {
	conf := stg.Smtp
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &smtpClient{
		settings: stg,
		dialer:   dialer,
	}
}

func (c *smtpClient) SendEmail(to, sub, msg string) error {
	slog.Info("Sending e-mail", "to", to, "subject", sub)
	conf := c.settings.Smtp
	mail := gomail.NewMessage()
	mail.SetHeader("From", conf.User)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", sub)
	mail.SetBody("text/plain", msg)

	if conf.From != "" {
		mail.SetHeader("From", conf.From)
	}

	if err := c.dialer.DialAndSend(mail); err != nil {
		return err
	}

	slog.Info("E-mail sent", "to", to, "subject", sub)

	return nil
}
