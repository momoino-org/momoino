package usermgt

import (
	"errors"
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// loginHandler is a struct that handles user login requests. It implements the
// core.HTTPRoute interface, providing methods to process login operations,
// including decoding the request, validating user credentials, and generating
// JWT tokens for authenticated users.
type renewSessionHandler struct {
	config         core.AppConfig
	logger         *slog.Logger
	db             *gorm.DB
	userService    UserService
	userRepository UserRepository
	jwtConfig      *core.JWTConfig
	sessionManager *scs.SessionManager
}

// LoginHandlerParams holds the parameters required to create a new loginHandler.
// It is used for dependency injection, allowing for easier testing and
// management of dependencies.
type RenewSessionHandlerParams struct {
	fx.In
	Logger         *slog.Logger
	Config         core.AppConfig
	SessionManager *scs.SessionManager
	DB             *gorm.DB
	UserService    UserService
	UserRepository UserRepository
}

// Ensure loginHandler implements the core.HTTPRoute interface.
var _ core.HTTPRoute = (*renewSessionHandler)(nil)

// NewLoginHandler initializes a new loginHandler with the provided parameters.
// It sets up the necessary dependencies for handling login requests.
//
// Parameters:
//   - params: A LoginHandlerParams struct containing dependencies needed by the login handler.
//
// Returns:
//   - A pointer to the newly created loginHandler.
func NewRenewSessionHandler(params LoginHandlerParams) *renewSessionHandler {
	handler := renewSessionHandler{
		config:         params.Config,
		logger:         params.Logger,
		sessionManager: params.SessionManager,
		db:             params.DB,
		userService:    params.UserService,
		userRepository: params.UserRepository,
		jwtConfig:      params.Config.GetJWTConfig(),
	}

	return &handler
}

func (h *renewSessionHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "POST /api/v1/token/renew",
		Wrappers: []func(http.Handler) http.Handler{
			h.sessionManager.LoadAndSave,
		},
	}
}

// ServeHTTP handles incoming HTTP requests for user login. It decodes the
// request body, validates the user credentials, generates JWT tokens, and
// returns appropriate responses based on the outcome of the operations.
//
// Parameters:
//   - w: An http.ResponseWriter to write the response to the client.
//   - r: An http.Request containing the client's request data.
func (h *renewSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	if !h.sessionManager.Exists(ctx, core.UIDKey) {
		h.userService.ClearAuthCookies(w, h.sessionManager)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgNeedToLogin).Build())

		return
	}

	userID := h.sessionManager.GetString(ctx, core.UIDKey)

	loggedUser, err := h.userRepository.FindUserByID(ctx, h.db, uuid.MustParse(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if destroySessionErr := h.sessionManager.Destroy(ctx); destroySessionErr != nil {
				h.logger.ErrorContext(ctx, "Cannot destroy session", core.DetailsLogAttr(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

				return
			}

			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgNeedToLogin).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by username", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	jwt, err := h.userService.GenerateJWT(h.sessionManager.Token(r.Context()), *loggedUser)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot create json web token", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	h.sessionManager.Put(r.Context(), core.UIDKey, loggedUser.ID.String())
	h.userService.SetAuthCookies(w, *jwt)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&LoginResponse{
		AccessToken: jwt.AccessToken.Value,
	}).Build())
}
