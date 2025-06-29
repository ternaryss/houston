package handlers

import (
	"net/http"
	"net/url"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views"
)

type DashboardHandler struct {
	webAppsStore types.WebAppsStore
}

func NewDashboardHandler(was types.WebAppsStore) *DashboardHandler {
	return &DashboardHandler{
		webAppsStore: was,
	}
}

func (h *DashboardHandler) Dashboard(res http.ResponseWriter, req *http.Request) {
	user := helpers.AuthPrincipal(req)
	params := url.Values{}
	params.Set("page", "1")
	params.Set("size", "10")
	page, err := cmd.NewGetWebAppsCmd(h.webAppsStore).Execute(params, user)

	if err != nil {
		helpers.InternalServerError(err, res, req)
		return
	}

	template := views.Dashboard(page, user)
	helpers.RenderPage(template, res, req)
}
