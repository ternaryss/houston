package handlers

import (
	"net/http"

	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/views"
)

type ErrorsHandler struct{}

func NewErrorsHandler() *ErrorsHandler {
	return &ErrorsHandler{}
}

func (h *ErrorsHandler) NotFoundError(res http.ResponseWriter, req *http.Request) {
	template := views.NotFoundError()
	helpers.RenderPage(template, res, req)
}

func (h *ErrorsHandler) InternalServerError(res http.ResponseWriter, req *http.Request) {
	template := views.InternalServerError()
	helpers.RenderPage(template, res, req)
}
