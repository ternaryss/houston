package helpers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func Render(cpt templ.Component, res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	if err := cpt.Render(req.Context(), res); err != nil {
		slog.Error("Rendering error", "err", err)
		os.Exit(1)
	}
}
