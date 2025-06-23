package helpers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/ternaryss/houston/views"
)

func RenderPage(cpt templ.Component, res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	view := views.Base(cpt)

	if err := view.Render(req.Context(), res); err != nil {
		slog.Error("Rendering page error", "err", err)
		os.Exit(1)
	}
}

func Render(cpt templ.Component, res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	if err := cpt.Render(req.Context(), res); err != nil {
		slog.Error("Rendering component error", "err", err)
		os.Exit(1)
	}
}
