package types

import "fmt"

type Message struct {
	Subject string
	Content string
}

func NewAliveMessage(app *WebApp, health *HealthCheck) Message {
	sub := fmt.Sprintf("%s health check notification", app.Name)
	msg := fmt.Sprintf("Web application %s is alive - status = %d, expected = %d.", app.Name, health.Status, app.Status)

	return Message{
		Subject: sub,
		Content: msg,
	}
}

func NewNotAliveMessage(app *WebApp, health *HealthCheck) Message {
	sub := fmt.Sprintf("%s health check notification", app.Name)
	msg := fmt.Sprintf("Web application %s is down - status = %d, expected = %d.", app.Name, health.Status, app.Status)

	return Message{
		Subject: sub,
		Content: msg,
	}
}
