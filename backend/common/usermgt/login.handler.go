package usermgt

import (
	"context"
	"crypto/rsa"
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
	logger         core.Logger
	db             *gorm.DB
	userRepository UserRepository
	jwtConfig      *core.JWTConfig
	jwtPrivateKey  *rsa.PrivateKey
}

type LoginHandlerParams struct {
	fx.In

	AppLifeCycle   fx.Lifecycle
	Config         core.AppConfig
	Logger         core.Logger
	DB             *gorm.DB
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

var (
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
	ErrCannotGenerateToken   = errors.New("cannot generate token")
)

func NewLoginHandler(params LoginHandlerParams) *loginHandler {
	handler := loginHandler{
		config:         params.Config,
		logger:         params.Logger,
		db:             params.DB,
		userRepository: params.UserRepository,
		jwtConfig:      params.Config.GetJWTConfig(),
	}

	params.AppLifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(handler.jwtConfig.PrivateKey)

			if err != nil {
				return err
			}

			handler.jwtPrivateKey = rsaPrivateKey

			return nil
		},
	})

	return &handler
}

func (h *loginHandler) Pattern() string {
	return "POST /api/v1/login"
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseBuilder := core.NewResponseBuilder()

	var requestBody LoginRequestBody

	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		h.logger.ErrorContext(ctx, "Something went wrong when trying to decode request body", slog.Any("error", err))
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, responseBuilder.MessageID("S0400").Build())

		return
	}

	loggedUser, err := h.userRepository.FindUserByEmail(ctx, h.db, requestBody.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, responseBuilder.MessageID("E-0001").Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by email", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID("U-0001").Build())

		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(loggedUser.Password), []byte(requestBody.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, responseBuilder.MessageID("E-0001").Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when comparing password", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID("U-0001").Build())

		return
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":       loggedUser.ID,
		"exp":       now.Add(time.Duration(h.jwtConfig.AccessTokenExpiresIn)).Unix(),
		"auth_time": now.Unix(),
	})

	tokenString, err := token.SignedString(h.jwtPrivateKey)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot sign access token", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID("U-0001").Build())

		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":       loggedUser.ID,
		"exp":       now.Add(time.Duration(h.jwtConfig.RefreshTokenExpiresIn)).Unix(),
		"auth_time": now.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(h.jwtPrivateKey)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot sign refresh token", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID("U-0001").Build())

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&LoginResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}).Build())
}
