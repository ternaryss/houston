package handlers

import (
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views"
	"github.com/ternaryss/houston/views/components"
)

type UsersHandler struct {
	usersStore types.UsersStore
}

func NewUsersHandler(urs types.UsersStore) *UsersHandler {
	return &UsersHandler{
		usersStore: urs,
	}
}

func (h *UsersHandler) SignUp(res http.ResponseWriter, req *http.Request) {
	// TODO: redirect to dashboard if Signed in
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

		// TODO: redirect to Sign in
		helpers.Redirect("/", res, req)
	}

	template := views.SignUp(form)
	helpers.RenderPage(template, res, req)
}
