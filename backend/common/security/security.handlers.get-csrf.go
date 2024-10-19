package security

import (
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/gorilla/csrf"
)

type getCsrfTokenHandler struct{}

var _ core.HTTPRoute = (*getCsrfTokenHandler)(nil)

func NewGetCsrfTokenHandler() *getCsrfTokenHandler {
	return &getCsrfTokenHandler{}
}

func (h *getCsrfTokenHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "GET /api/v1/csrf-token",
	}
}

func (h *getCsrfTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Csrf-Token", csrf.Token(r))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, core.NewResponseBuilder(r).Build())
}
