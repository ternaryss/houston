package types

import "time"

const (
	Interval5M  string = "5M"
	Interval15M string = "15M"
	Interval1H  string = "1H"
)

var Intervals = []string{Interval5M, Interval15M, Interval1H}

type User struct {
	Email      string
	Password   string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewUser(email, password string) *User {
	now := time.Now().UTC()

	return &User{
		Email:      email,
		Password:   password,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

type WebApp struct {
	Id         string
	Name       string
	Url        string
	Status     int
	Interval   string
	UserEmail  string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewWebApp(name, url, interval, userEmail string, status int) *WebApp {
	now := time.Now().UTC()

	return &WebApp{
		Name:       name,
		Url:        url,
		Status:     status,
		Interval:   interval,
		UserEmail:  userEmail,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

type Subscriber struct {
	Id         int64
	WebAppId   string
	Email      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewSubscriber(webAppId, email string) *Subscriber {
	now := time.Now().UTC()

	return &Subscriber{
		WebAppId:   webAppId,
		Email:      email,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}
