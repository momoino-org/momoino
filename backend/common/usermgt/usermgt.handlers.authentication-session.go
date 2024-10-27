package usermgt

import (
	"net/http"
	"wano-island/common/core"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/render"
	"go.uber.org/fx"
)

type createLoginSessionHandler struct {
	sessionManager      *scs.SessionManager
	loginSessionManager *scs.SessionManager
}

type CreateLoginSessionHandlerParams struct {
	fx.In
	SessionManager           *scs.SessionManager
	ShortLivedSessionManager *scs.SessionManager `name:"loginSessionManager"`
}

func NewCreateLoginSessionHandler(params CreateLoginSessionHandlerParams) *createLoginSessionHandler {
	return &createLoginSessionHandler{
		sessionManager:      params.SessionManager,
		loginSessionManager: params.ShortLivedSessionManager,
	}
}

func (h *createLoginSessionHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "POST /api/v1/authentication/session",
		Wrappers: []func(http.Handler) http.Handler{
			h.loginSessionManager.LoadAndSave,
			h.sessionManager.LoadAndSave,
		},
	}
}

func (h *createLoginSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	// Clean current session to start fresh session
	if destroySessionErr := h.loginSessionManager.RenewToken(reqCtx); destroySessionErr != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgCannotDestroySession).Build())

		return
	}

	// Clean current session to start fresh session
	if destroySessionErr := h.sessionManager.Destroy(reqCtx); destroySessionErr != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgCannotDestroySession).Build())

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   core.SessionCookie,
		MaxAge: -1,
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Build())
}
