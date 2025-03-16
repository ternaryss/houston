package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type middleware func(http.Handler) http.HandlerFunc

func middlewaresChain(mid ...middleware) middleware {
	return func(nxt http.Handler) http.HandlerFunc {
		for i := len(mid) - 1; i >= 0; i-- {
			nxt = mid[i](nxt)
		}

		return nxt.ServeHTTP
	}
}

func requestTimeMiddleware(nxt http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		nxt.ServeHTTP(res, req)
		duration := time.Since(start)
		slog.Info("Request info", "method", req.Method, "path", req.RequestURI, "time", fmt.Sprintf("%dms", duration.Microseconds()))
	}
}
