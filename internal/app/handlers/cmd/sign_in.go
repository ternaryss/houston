package cmd

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
	"golang.org/x/crypto/bcrypt"
)

type signInCmd struct {
	settings   settings.Settings
	usersStore types.UsersStore
}

func NewSignInCmd(stg settings.Settings, urs types.UsersStore) *signInCmd {
	return &signInCmd{
		settings:   stg,
		usersStore: urs,
	}
}

func (c *signInCmd) Execute(frm types.SignInForm) (string, error) {
	slog.Info("Signing user in", "email", frm.Email)
	frm.Validate()

	if len(frm.Errors) > 0 {
		slog.Info("Sign in form is not valid", "errors", frm.Errors)
		return "", nil
	}

	user, err := c.usersStore.GetByEmail(frm.Email, nil)

	if err != nil {
		if err == sql.ErrNoRows {
			frm.Errors["password"] = types.NewFieldError("password", "Invalid address e-mail or password.")
			slog.Info("Sign in form is not valid", "errors", frm.Errors)
			return "", nil
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(frm.Password)); err != nil {
		frm.Errors["password"] = types.NewFieldError("password", "Invalid address e-mail or password.")
		slog.Info("Sign in form is not valid", "errors", frm.Errors)
		return "", nil
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(c.settings.Authorization.ExpiresAfter) * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   user.Email,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.settings.Authorization.Secret))

	if err != nil {
		return "", err
	}

	slog.Info("User signed in", "email", user.Email, "expiresAt", expiresAt)

	return signedToken, nil
}
