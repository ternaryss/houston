package helpers

import (
	"net/http"

	"github.com/ternaryss/houston/internal/app/types"
)

func AuthPrincipal(req *http.Request) string {
	userCtx := req.Context().Value(types.CtxUserKey)

	if user, ok := userCtx.(string); ok {
		return user
	}

	return ""
}
