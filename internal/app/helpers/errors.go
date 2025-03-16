package helpers

import (
	"log/slog"
	"net/http"
)

func NotFoundError(res http.ResponseWriter, req *http.Request) {
	slog.Warn("Path not found", "path", req.URL.Path)
	Redirect("/not-found", res, req)
}

func InternalServerError(err error, res http.ResponseWriter, req *http.Request) {
	slog.Error("Handling unexpected error", "err", err)
	Redirect("/error", res, req)
}
