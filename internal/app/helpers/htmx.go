package helpers

import "net/http"

func IsHtmxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}

func Redirect(loc string, res http.ResponseWriter, req *http.Request) {
	if IsHtmxRequest(req) {
		res.Header().Set("HX-Redirect", loc)
		res.WriteHeader(http.StatusSeeOther)
		return
	}

	res.Header().Set("Location", loc)
	res.WriteHeader(http.StatusSeeOther)
}
