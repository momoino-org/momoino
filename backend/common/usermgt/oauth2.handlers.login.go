//nolint:dupl // This is different with OAuth2CallbackHandler
package usermgt

import (
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"go.uber.org/fx"
)

// oauth2LoginHandler handles OAuth2 login requests by routing them to the appropriate provider.
// It maps provider names to OAuth2Provider implementations and delegates the authorization process.
type oauth2LoginHandler struct {
	providers map[string]OAuth2Provider
}

// OAuth2LoginHandlerParams defines the dependencies required to initialize an oauth2LoginHandler.
type OAuth2LoginHandlerParams struct {
	fx.In
	GoogleProvider OAuth2Provider `name:"google_provider"`
}

// Ensure that oauth2LoginHandler implements the core.HTTPRoute interface.
// This guarantees that oauth2LoginHandler conforms to the expected route handler structure.
var _ core.HTTPRoute = (*oauth2LoginHandler)(nil)

// NewOAuth2LoginHandler creates and returns a new instance of oauth2LoginHandler.
// It initializes the handler with a map of available OAuth2 providers.
//
// Params:
//   - params: The dependencies required for the handler, including OAuth2 providers.
//
// Returns:
//   - A new oauth2LoginHandler instance.
func NewOAuth2LoginHandler(params OAuth2LoginHandlerParams) *oauth2LoginHandler {
	return &oauth2LoginHandler{
		providers: map[string]OAuth2Provider{
			GoogleProviderName: params.GoogleProvider,
		},
	}
}

func (h *oauth2LoginHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "GET /api/v1/login/providers/{provider}",
	}
}

// ServeHTTP processes the incoming HTTP request and delegates authorization to the appropriate OAuth2 provider.
// It looks up the provider based on the {provider} path parameter. If the provider is found, the request is
// passed to the provider's AuthorizeHandler; otherwise, a 404 Not Found response is returned.
//
// Params:
//   - w: The HTTP response writer to send the response.
//   - r: The HTTP request containing the provider name in the URL.
//
// Behavior:
//   - If the requested provider exists in the handler's providers map, it invokes the provider's AuthorizeHandler.
//   - If the provider is not found, it returns a 404 Not Found error with a response message.
func (h *oauth2LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if provider, exists := h.providers[r.PathValue("provider")]; exists {
		provider.AuthorizeHandler(w, r)
		return
	}

	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgRouteNotFound).Build())
}
