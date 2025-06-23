package handlers

import (
	"net/http"
	"time"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views"
	"github.com/ternaryss/houston/views/components"
)

type UsersHandler struct {
	settings   settings.Settings
	usersStore types.UsersStore
}

func NewUsersHandler(stg settings.Settings, urs types.UsersStore) *UsersHandler {
	return &UsersHandler{
		settings:   stg,
		usersStore: urs,
	}
}

func (h *UsersHandler) SignIn(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)

	if user != "" {
		helpers.Redirect("/", res, req)
		return
	}

	signUpEnabled := h.settings.Authorization.SignUp.Enabled
	form, err := types.NewSignInForm(nil)

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	if req.Method == http.MethodPost {
		form, err = types.NewSignInForm(req)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		token, err := cmd.NewSignInCmd(h.settings, h.usersStore).Execute(form)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		if len(form.Errors) > 0 {
			helpers.Render(components.SignInForm(form, signUpEnabled), res, req)
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		http.SetCookie(res, cookie)
		helpers.Redirect("/", res, req)
		return
	}

	template := views.SignIn(form, signUpEnabled)
	helpers.RenderPage(template, res, req)
}

func (h *UsersHandler) SignOut(res http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(res, cookie)
	helpers.Redirect("/sign-in", res, req)
}

func (h *UsersHandler) SignUp(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)

	if user != "" {
		helpers.Redirect("/", res, req)
		return
	}

	form, err := types.NewSignUpForm(nil)

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	if req.Method == http.MethodPost {
		form, err = types.NewSignUpForm(req)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		if err := cmd.NewSignUpCmd(h.usersStore).Execute(form); err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		if len(form.Errors) > 0 {
			helpers.Render(components.SignUpForm(form), res, req)
			return
		}

		helpers.Redirect("/sign-in", res, req)
		return
	}

	template := views.SignUp(form)
	helpers.RenderPage(template, res, req)
}
