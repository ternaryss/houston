package types

import "time"

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
	UserEmail  string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewWebApp(name, url, userEmail string) *WebApp {
	now := time.Now().UTC()

	return &WebApp{
		Name:       name,
		Url:        url,
		UserEmail:  userEmail,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}
