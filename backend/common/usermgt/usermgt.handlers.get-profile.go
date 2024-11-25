package usermgt

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// profileHandler handles requests related to user profiles.
// It implements the core.HTTPRoute interface to define the HTTP route for retrieving user profile information.
type profileHandler struct {
	logger *slog.Logger
	db     *gorm.DB
}

// ProfileHandlerParams contains the dependencies needed to create a new profileHandler.
type ProfileHandlerParams struct {
	fx.In
	Logger *slog.Logger
	DB     *gorm.DB
}

// ProfileResponse represents the structure of the response containing user profile information.
type ProfileResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Locale    string `json:"locale"`
}

// Ensure profileHandler implements the core.HTTPRoute interface.
var _ core.HTTPRoute = (*profileHandler)(nil)

// NewProfileHandler creates a new instance of profileHandler with the provided parameters.
//
// Params:
//   - params: Dependencies required by the handler.
//
// Returns:
//   - *profileHandler: A new instance of profileHandler.
func NewProfileHandler(params ProfileHandlerParams) *profileHandler {
	return &profileHandler{
		logger: params.Logger,
		db:     params.DB,
	}
}

func (h *profileHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern:   "GET /api/v1/profile",
		IsPrivate: true,
	}
}

// ServeHTTP handles incoming HTTP requests for retrieving user profile information.
//
// Params:
//   - w: The http.ResponseWriter to send the response.
//   - r: The incoming http.Request containing the request for user profile.
func (h *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseBuilder := core.NewResponseBuilder(r)

	authUser := core.MustGetAuthUserFromRequest(r)

	// Respond with the user's profile information.
	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&ProfileResponse{
		ID:        authUser.GetID().String(),
		Username:  authUser.GetUsername(),
		Email:     authUser.GetEmail(),
		FirstName: authUser.GetGivenName(),
		LastName:  authUser.GetFamilyName(),
		Locale:    authUser.GetLocale(),
	}).Build())
}
