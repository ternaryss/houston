package handlers

import (
	"net/http"
	"strconv"

	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/views"
	"github.com/ternaryss/houston/views/components"
)

func Hello(res http.ResponseWriter, req *http.Request) {
	count := 0
	template := views.Index(count)

	if helpers.IsHtmxRequest(req) {
		qry := req.URL.Query()
		count, _ = strconv.Atoi(qry.Get("count"))
		count++
		template = components.Counter(count)
	}

	helpers.Render(template, res, req)
}
