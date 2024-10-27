package usermgt

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/alexedwards/scs/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines an interface for user authentication and management operations.
// It includes methods for comparing passwords, hashing passwords, generating JSON Web Tokens (JWT),
// and setting authentication cookies.
type UserService interface {
	// ComparePassword checks if the provided plain password matches the hashed password.
	//
	// Params:
	//   - ctx: The context for managing request-scoped values and cancellation signals.
	//   - password: The plain text password to compare.
	//   - hashedPassword: The hashed password to compare against.
	//
	// Returns:
	//   - error: Returns nil if the passwords match, or an error if they do not.
	ComparePassword(ctx context.Context, password []byte, hasedPassword []byte) error

	// HashPassword generates a hashed version of the provided password.
	//
	// Params:
	//   - ctx: The context for managing request-scoped values and cancellation signals.
	//   - password: The plain text password to be hashed.
	//
	// Returns:
	//   - []byte: The hashed password as a byte slice, or nil if hashing fails.
	//   - error: Returns nil if hashing is successful, or an error if it fails.
	HashPassword(ctx context.Context, password string) ([]byte, error)

	// GenerateJWT creates a JSON Web Token for the specified user.
	//
	// Params:
	//   - user: The UserModel containing user details used to generate the token.
	//
	// Returns:
	//   - *JWT: A pointer to the generated JWT, or nil if generation fails.
	//   - error: Returns nil if JWT generation is successful, or an error if it fails.
	GenerateJWT(sessionID string, user UserModel) (*JWT, error)

	// SetAuthCookies establishes authentication cookies in the HTTP response.
	//
	// Params:
	//   - w: The http.ResponseWriter used to set cookies.
	//   - jwt: The JSON Web Token to be set as a cookie.
	SetAuthCookies(w http.ResponseWriter, jwt JWT)

	ClearAuthCookies(w http.ResponseWriter, sessionManager *scs.SessionManager)
}

// userService is an implementation of the UserService interface.
type userService struct {
	logger    *slog.Logger
	config    core.AppConfig
	jwtConfig *core.JWTConfig
}

// UserServiceParams defines the dependencies required to create a UserService.
type UserServiceParams struct {
	fx.In
	Logger *slog.Logger
	Config core.AppConfig
}

// Token represents a JWT token, which contains its value and expiration time.
type Token struct {
	// The string representation of the token.
	Value string

	// The expiration time of the token.
	ExpiredAt time.Time
}

// JWT represents a structure containing both the access token and refresh token used for authentication.
type JWT struct {
	// The access token used for short-term authentication.
	AccessToken Token

	// The refresh token used to generate a new access token after expiration.
	RefreshToken Token
}

// _ ensures that userService implements the UserService interface.
// This is a compile-time check to confirm that userService provides all required methods.
var _ UserService = (*userService)(nil)

// NewUserService creates and returns a new instance of userService.
// It initializes the service with the provided logger and JWT configuration.
func NewUserService(params UserServiceParams) *userService {
	return &userService{
		logger:    params.Logger,
		config:    params.Config,
		jwtConfig: params.Config.GetJWTConfig(),
	}
}

func (s *userService) ComparePassword(ctx context.Context, password []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	if err == nil {
		return nil
	}

	if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		s.logger.ErrorContext(ctx, "Something went wrong when comparing password", core.DetailsLogAttr(err))
		return err
	}

	return err
}

func (s *userService) HashPassword(ctx context.Context, password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		s.logger.ErrorContext(ctx, "Something went wrong when hashing password", core.DetailsLogAttr(err))
		return nil, err
	}

	return hashedPassword, nil
}

func (s *userService) generateAccessToken(sessionID string, user UserModel, nowFn func() time.Time) (*Token, error) {
	now := nowFn()
	claims := core.JWTCustomClaims{
		SessionID:         sessionID,
		Email:             user.Email,
		PreferredUsername: user.Username,
		GivenName:         user.FirstName,
		FamilyName:        user.LastName,
		Locale:            user.Locale,
		Roles:             []string{},
		Permissions:       []string{},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.AccessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(s.jwtConfig.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     tokenString,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *userService) generateRefreshToken(sessionID string, user UserModel, nowFn func() time.Time) (*Token, error) {
	now := nowFn()
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.RefreshTokenExpiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	refreshTokenString, err := refreshToken.SignedString(s.jwtConfig.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     refreshTokenString,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *userService) GenerateJWT(sessionID string, user UserModel) (*JWT, error) {
	now := time.Now()

	accessToken, signedAccessTokenErr := s.generateAccessToken(sessionID, user, func() time.Time {
		return now
	})

	if signedAccessTokenErr != nil {
		return nil, signedAccessTokenErr
	}

	refreshToken, signedRefreshTokenErr := s.generateRefreshToken(sessionID, user, func() time.Time {
		return now
	})

	if signedRefreshTokenErr != nil {
		return nil, signedRefreshTokenErr
	}

	return &JWT{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s *userService) SetAuthCookies(w http.ResponseWriter, jwt JWT) {
	http.SetCookie(w, &http.Cookie{
		Name:     core.IdentityCookie,
		Value:    jwt.AccessToken.Value,
		HttpOnly: true,
		Domain:   s.config.GetHost(),
		Secure:   s.config.IsHTTPS(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  jwt.AccessToken.ExpiredAt,
	})
}

func (s *userService) ClearAuthCookies(w http.ResponseWriter, sessionManager *scs.SessionManager) {
	http.SetCookie(w, &http.Cookie{
		Name:     core.IdentityCookie,
		HttpOnly: true,
		Domain:   s.config.GetHost(),
		Secure:   s.config.IsHTTPS(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     sessionManager.Cookie.Name,
		HttpOnly: sessionManager.Cookie.HttpOnly,
		SameSite: sessionManager.Cookie.SameSite,
		Path:     sessionManager.Cookie.Path,
		Domain:   sessionManager.Cookie.Domain,
		Secure:   sessionManager.Cookie.Secure,
		MaxAge:   -1,
	})
}
