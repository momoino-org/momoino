package usermgt

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginHandler struct {
	config         core.AppConfig
	logger         *slog.Logger
	db             *gorm.DB
	userService    UserService
	userRepository UserRepository
	jwtConfig      *core.JWTConfig
}

type LoginHandlerParams struct {
	fx.In

	Logger         *slog.Logger
	Config         core.AppConfig
	DB             *gorm.DB
	UserService    UserService
	UserRepository UserRepository
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewLoginHandler(params LoginHandlerParams) *loginHandler {
	handler := loginHandler{
		config:         params.Config,
		logger:         params.Logger,
		db:             params.DB,
		userService:    params.UserService,
		userRepository: params.UserRepository,
		jwtConfig:      params.Config.GetJWTConfig(),
	}

	return &handler
}

func (h *loginHandler) Pattern() string {
	return "POST /api/v1/login"
}

func (h *loginHandler) IsPrivateRoute() bool {
	return false
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	var requestBody LoginRequestBody
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		h.logger.ErrorContext(ctx, "Something went wrong when trying to decode request body", slog.Any("details", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	loggedUser, err := h.userRepository.FindUserByEmail(ctx, h.db, requestBody.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by email", slog.Any("details", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	if err = h.userService.ComparePassword(
		r.Context(),
		[]byte(requestBody.Password),
		[]byte(loggedUser.Password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, core.JWTCustomClaims{
		Email:             loggedUser.Email,
		PreferredUsername: loggedUser.Username,
		GivenName:         loggedUser.FirstName,
		FamilyName:        loggedUser.LastName,
		Locale:            loggedUser.Locale,
		Roles:             []string{},
		Permissions:       []string{},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   loggedUser.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(h.jwtConfig.AccessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})

	tokenString, err := token.SignedString(h.jwtConfig.PrivateKey)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot sign access token", slog.Any("details", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   loggedUser.ID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(h.jwtConfig.AccessTokenExpiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	})

	refreshTokenString, err := refreshToken.SignedString(h.jwtConfig.PrivateKey)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot sign refresh token", slog.Any("details", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&LoginResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}).Build())
}
