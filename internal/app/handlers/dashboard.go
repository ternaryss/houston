package handlers

import (
	"net/http"

	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/views"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Dashboard(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(types.CtxUserKey).(string)
	template := views.Dashboard(user)
	helpers.RenderPage(template, res, req)
}
