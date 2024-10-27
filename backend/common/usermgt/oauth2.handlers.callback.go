//nolint:dupl // This is different with OAuth2LoginHandler
package usermgt

import (
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"go.uber.org/fx"
)

// oauth2LoginCallbackHandler processes OAuth2 callback requests for login.
// It routes the callback requests to the appropriate OAuth2 provider's callback handler based on the provider name.
type oauth2LoginCallbackHandler struct {
	providers map[string]OAuth2Provider
}

// OAuth2LoginCallbackHandlerParams defines the dependencies required to initialize an oauth2LoginCallbackHandler.
type OAuth2LoginCallbackHandlerParams struct {
	fx.In
	GoogleProvider OAuth2Provider `name:"google_provider"`
}

// Ensure that oauth2LoginCallbackHandler implements the core.HTTPRoute interface.
// This guarantees that oauth2LoginCallbackHandler conforms to the expected route handler structure.
var _ core.HTTPRoute = (*oauth2LoginCallbackHandler)(nil)

// NewOAuth2LoginCallbackHandler creates and returns a new instance of oauth2LoginCallbackHandler.
// It initializes the handler with a map of available OAuth2 providers.
//
// Params:
//   - params: The dependencies required for the handler, including OAuth2 providers.
//
// Returns:
//   - A new oauth2LoginCallbackHandler instance.
func NewOAuth2LoginCallbackHandler(params OAuth2LoginCallbackHandlerParams) *oauth2LoginCallbackHandler {
	return &oauth2LoginCallbackHandler{
		providers: map[string]OAuth2Provider{
			GoogleProviderName: params.GoogleProvider,
		},
	}
}

func (h *oauth2LoginCallbackHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "GET /api/v1/login/providers/{provider}/callback",
	}
}

// ServeHTTP processes the incoming HTTP callback request and delegates the handling
// to the appropriate OAuth2 provider's callback handler. It looks up the provider
// based on the {provider} path parameter. If the provider is found, the request is
// passed to the provider's CallbackHandler; otherwise, a 404 Not Found response is returned.
//
// Params:
//   - w: The HTTP response writer to send the response.
//   - r: The HTTP request containing the provider name in the URL.
//
// Behavior:
//   - If the requested provider exists in the handler's providers map, it invokes the provider's CallbackHandler.
//   - If the provider is not found, it returns a 404 Not Found error with a response message.
func (h *oauth2LoginCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if provider, exists := h.providers[r.PathValue("provider")]; exists {
		provider.CallbackHandler(w, r)
		return
	}

	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgRouteNotFound).Build())
}
