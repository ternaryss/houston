package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views"
	"github.com/ternaryss/houston/views/components"
)

type WebAppsHandler struct {
	webAppsStore types.WebAppsStore
}

func NewWebAppsHandler(was types.WebAppsStore) *WebAppsHandler {
	return &WebAppsHandler{
		webAppsStore: was,
	}
}

func (h *WebAppsHandler) AddWebApp(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	form, err := types.NewWebAppForm(nil)

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	if req.Method == http.MethodPost {
		form, err = types.NewWebAppForm(req)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		id, err := cmd.NewAddWebAppCmd(h.webAppsStore).Execute(form, user)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		if len(form.Errors) > 0 {
			helpers.Render(components.WebAppForm(form, types.AddMode), res, req)
			return
		}

		helpers.Redirect(fmt.Sprintf("/web-apps/%s", id), res, req)
		return
	}

	template := views.WebApp(form, nil, user, types.AddMode)
	helpers.RenderPage(template, res, req)
}

func (h *WebAppsHandler) GetWebApps(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	page, err := cmd.NewGetWebAppsCmd(h.webAppsStore).Execute(req.URL.Query(), user)

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	template := components.WebAppsList(page)
	helpers.Render(template, res, req)
}

func (h *WebAppsHandler) GetWebApp(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	id := req.PathValue("id")
	app, err := cmd.NewGetWebAppCmd(h.webAppsStore).Execute(id, user)

	if err != nil {
		if err == sql.ErrNoRows {
			helpers.NotFoundError(res, req)
			return
		}

		helpers.InternalServerError(err, res, req)
		return
	}

	template := views.WebApp(types.WebAppForm{}, app, user, types.ReadMode)
	helpers.RenderPage(template, res, req)
}

func (h *WebAppsHandler) EditWebApp(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	id := req.PathValue("id")

	if req.Method == http.MethodPut {
		form, err := types.NewWebAppForm(req)

		if err != nil {
			helpers.InternalServerError(err, res, req)
			return
		}

		if err := cmd.NewEditWebAppCmd(h.webAppsStore).Execute(form, id, user); err != nil {
			if err == sql.ErrNoRows {
				helpers.NotFoundError(res, req)
				return
			}

			helpers.InternalServerError(err, res, req)
			return
		}

		if len(form.Errors) > 0 {
			form.Id = id
			helpers.Render(components.WebAppForm(form, types.EditMode), res, req)
			return
		}

		helpers.Redirect(fmt.Sprintf("/web-apps/%s", id), res, req)
		return
	}

	app, err := cmd.NewGetWebAppCmd(h.webAppsStore).Execute(id, user)

	if err != nil {
		if err == sql.ErrNoRows {
			helpers.NotFoundError(res, req)
			return
		}

		helpers.InternalServerError(err, res, req)
		return
	}

	form := types.WebAppForm{
		Id:     app.Id,
		Name:   app.Name,
		Url:    app.Url,
		Errors: make(map[string]types.FieldError),
	}
	template := views.WebApp(form, nil, user, types.EditMode)
	helpers.RenderPage(template, res, req)
}

func (h *WebAppsHandler) DeleteWebApp(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	id := req.PathValue("id")

	if err := cmd.NewDeleteWebAppCmd(h.webAppsStore).Execute(id, user); err != nil {
		if err == sql.ErrNoRows {
			helpers.NotFoundError(res, req)
			return
		}

		helpers.InternalServerError(err, res, req)
		return
	}

	helpers.Redirect("/", res, req)
}
