package types

type EmailClient interface {
	SendEmail(to, sub, msg string) error
}
