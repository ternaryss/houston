package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ternaryss/houston/internal/app/helpers"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
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

func authorizationMiddleware(stg settings.Settings) middleware {
	paths := []string{"/not-found", "/error", "/sign-in", "/sign-up"}
	resources := []string{"/public"}

	return func(nxt http.Handler) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			user := ""

			if cookie, err := req.Cookie("token"); err == nil {
				claims := &jwt.RegisteredClaims{}
				token, err := jwt.ParseWithClaims(
					cookie.Value,
					claims,
					func(tkn *jwt.Token) (any, error) {
						return []byte(stg.Authorization.Secret), nil
					},
				)

				if err == nil && token.Valid {
					user = claims.Subject
				}
			}

			ctx := req.Context()

			if user != "" {
				ctx = context.WithValue(ctx, types.CtxUserKey, user)
			}

			for _, resource := range resources {
				if strings.HasPrefix(req.URL.Path, resource) {
					nxt.ServeHTTP(res, req.WithContext(ctx))
					return
				}
			}

			if slices.Contains(paths, req.URL.Path) {
				nxt.ServeHTTP(res, req.WithContext(ctx))
				return
			}

			if user == "" {
				helpers.Redirect("/sign-in", res, req)
				return
			}

			nxt.ServeHTTP(res, req.WithContext(ctx))
		}
	}
}
