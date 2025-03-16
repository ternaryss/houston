package types

import "time"

type User struct {
	Email      string
	Password   string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewUser(email, password string) *User {
	now := time.Now()

	return &User{
		Email:      email,
		Password:   password,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}
