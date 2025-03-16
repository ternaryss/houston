package cmd

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/ternaryss/houston/internal/app/types"
	"golang.org/x/crypto/bcrypt"
)

type signUpCmd struct {
	usersStore types.UsersStore
}

func NewSignUpCmd(urs types.UsersStore) *signUpCmd {
	return &signUpCmd{
		usersStore: urs,
	}
}

func (c *signUpCmd) Execute(frm types.SignUpFrom) error {
	slog.Info("Signing up new user", "email", frm.Email)
	frm.Validate()

	if len(frm.Errors) > 0 {
		slog.Info("Sign up form is not valid", "errors", frm.Errors)
		return nil
	}

	email := strings.ToLower(strings.TrimSpace(frm.Email))
	password := strings.TrimSpace(frm.Password)
	exists, err := c.usersStore.GetByEmail(email)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if exists != nil {
		frm.Errors["email"] = types.NewFieldError("email", "Already in use.")
		slog.Info("Sign up form is not valid", "errors", frm.Errors)
		return nil
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := types.NewUser(email, string(passHash))

	if _, err := c.usersStore.Insert(user); err != nil {
		return err
	}

	slog.Info("New user signed up", "email", user.Email)

	return nil
}
