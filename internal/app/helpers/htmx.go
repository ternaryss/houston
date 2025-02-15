package helpers

import "net/http"

func IsHtmxRequest(req *http.Request) bool {
	return req.Header.Get("HX-Request") == "true"
}
