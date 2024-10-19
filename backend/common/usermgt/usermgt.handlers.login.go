package usermgt

import (
	"errors"
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/render"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// loginHandler is a struct that handles user login requests. It implements the
// core.HTTPRoute interface, providing methods to process login operations,
// including decoding the request, validating user credentials, and generating
// JWT tokens for authenticated users.
type loginHandler struct {
	config                   core.AppConfig
	logger                   *slog.Logger
	db                       *gorm.DB
	userService              UserService
	userRepository           UserRepository
	jwtConfig                *core.JWTConfig
	sessionManager           *scs.SessionManager
	shortLivedSessionManager *scs.SessionManager
}

// LoginHandlerParams holds the parameters required to create a new loginHandler.
// It is used for dependency injection, allowing for easier testing and
// management of dependencies.
type LoginHandlerParams struct {
	fx.In
	Logger                   *slog.Logger
	Config                   core.AppConfig
	SessionManager           *scs.SessionManager
	ShortLivedSessionManager *scs.SessionManager `name:"shortLivedSessionManager"`
	DB                       *gorm.DB
	UserService              UserService
	UserRepository           UserRepository
}

// LoginRequestBody defines the structure of the request body for login requests.
// It contains the username and password submitted by the user.
type LoginRequestBody struct {
	//  The username of the user attempting to log in.
	Username string `json:"username"`
	// The password of the user attempting to log in.
	Password string `json:"password"`
}

// LoginResponse represents the structure of the response sent back to the client
// after a successful login. It contains the access and refresh tokens.

type LoginResponse struct {
	// The JWT access token for authorizing future requests.
	AccessToken string `json:"accessToken"`
}

// Ensure loginHandler implements the core.HTTPRoute interface.
var _ core.HTTPRoute = (*loginHandler)(nil)

// NewLoginHandler initializes a new loginHandler with the provided parameters.
// It sets up the necessary dependencies for handling login requests.
//
// Parameters:
//   - params: A LoginHandlerParams struct containing dependencies needed by the login handler.
//
// Returns:
//   - A pointer to the newly created loginHandler.
func NewLoginHandler(params LoginHandlerParams) *loginHandler {
	handler := loginHandler{
		config:                   params.Config,
		logger:                   params.Logger,
		sessionManager:           params.SessionManager,
		shortLivedSessionManager: params.ShortLivedSessionManager,
		db:                       params.DB,
		userService:              params.UserService,
		userRepository:           params.UserRepository,
		jwtConfig:                params.Config.GetJWTConfig(),
	}

	return &handler
}

func (h *loginHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "POST /api/v1/login",
		Wrappers: []func(http.Handler) http.Handler{
			// h.shortLivedSessionManager.LoadAndSave,
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
func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	shortLivedSessionCookie, getShortLivedSessionCookieErr := r.Cookie(core.LoginSessionCookie)
	if getShortLivedSessionCookieErr != nil {
		h.logger.ErrorContext(ctx,
			"Cannot get short lived session cookie",
			core.DetailsLogAttr(getShortLivedSessionCookieErr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidSession).Build())

		return
	}

	shortLivedSessionCtx, loadShortLivedSessionErr := h.shortLivedSessionManager.Load(ctx, shortLivedSessionCookie.Value)
	if loadShortLivedSessionErr != nil {
		h.logger.ErrorContext(ctx, "Cannot load short lived session", core.DetailsLogAttr(loadShortLivedSessionErr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidSession).Build())

		return
	}

	var requestBody LoginRequestBody
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		h.logger.ErrorContext(ctx, "Something went wrong when trying to decode request body", core.DetailsLogAttr(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	loggedUser, err := h.userRepository.FindUserByUsername(ctx, h.db, requestBody.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by username", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	if err = h.userService.ComparePassword(
		r.Context(),
		[]byte(requestBody.Password),
		[]byte(*loggedUser.Password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when comparing password", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	if destroyShortSessionErr := h.shortLivedSessionManager.Destroy(shortLivedSessionCtx); destroyShortSessionErr != nil {
		h.logger.ErrorContext(ctx, "Cannot destroy short live session", core.DetailsLogAttr(destroyShortSessionErr))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgCannotDestroySession).Build())

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   core.LoginSessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if renewTokenErr := h.sessionManager.RenewToken(ctx); renewTokenErr != nil {
		h.logger.ErrorContext(ctx, "Cannot renew session token", core.DetailsLogAttr(renewTokenErr))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	jwt, err := h.userService.GenerateJWT(h.sessionManager.Token(ctx), *loggedUser)
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
