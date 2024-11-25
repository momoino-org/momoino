package httpsrv

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"wano-island/common/core"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

// extractAccessTokenFromRequest extracts the access token from the HTTP request header.
func extractAccessTokenFromRequest(req *http.Request) string {
	// Extract the access token from the request header.
	accessToken := strings.TrimPrefix(req.Header.Get(core.AuthorizationHeader), "Bearer ")

	return strings.TrimSpace(accessToken)
}

// IsBypassMiddleware checks if the JWT middleware should bypass authentication for the given HTTP request.
// It checks if the application is in testing mode and if a test user is present in the request.
//
// Parameters:
//   - r: The HTTP request to be processed.
//   - config: The application configuration object.
//
// Returns:
//   - A boolean value indicating whether the JWT middleware should bypass authentication.
//     Returns true if the application is in testing mode and a test user is present in the request.
//     Returns false otherwise.
func IsBypassMiddleware(r *http.Request, config core.AppConfig) bool {
	_, err := core.GetTestAuthUserFromRequest(r)

	return config.IsTesting() && !errors.Is(err, core.ErrAuthUserNotFound)
}

// WithJwtMiddleware is a middleware function that handles JWT authentication for HTTP requests.
// It extracts the access token from the request header, verifies it, and converts the claims into an AuthUser object.
// If the access token is missing or invalid, it returns an appropriate HTTP response.
// Otherwise, it sets the AuthUser object in the request context and calls the next handler.
func WithJwtMiddleware(
	bundle *i18n.Bundle,
	config core.AppConfig,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var authUser core.PrincipalUser

			//nolint:nestif // Will fix later.
			if IsBypassMiddleware(r, config) {
				authUser, _ = core.GetTestAuthUserFromRequest(r)
			} else {
				accessToken := extractAccessTokenFromRequest(r)

				// If the access token is missing, return an unauthorized response.
				if accessToken == "" {
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgNeedToLogin).Build())

					return
				}

				// Parse the access token without verifying it to extract the JWT claims.
				jwtMapClaims := core.JWTCustomClaims{}
				jwtParser := jwt.NewParser()
				k, err := keyfunc.NewDefaultCtx(
					r.Context(),
					[]string{config.GetKeycloakProvider().JwksURI},
				)

				if err != nil {
					logger.ErrorContext(r.Context(), "Cannot create keyfunc", core.DetailsLogAttr(err))
					render.Status(r, http.StatusInternalServerError)
					render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

					return
				}

				if _, err := jwtParser.ParseWithClaims(accessToken, &jwtMapClaims, k.KeyfuncCtx(r.Context())); err != nil {
					logger.ErrorContext(r.Context(), "Cannot parse jwt", core.DetailsLogAttr(err))
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

					return
				}

				// Convert the JWT claims into an AuthUser object.
				authUser = &core.AuthenticatedUser{
					ID:          jwtMapClaims.Subject,
					Username:    jwtMapClaims.PreferredUsername,
					Email:       jwtMapClaims.Email,
					FamilyName:  jwtMapClaims.FamilyName,
					GivenName:   jwtMapClaims.GivenName,
					Locale:      lo.Ternary(lo.IsEmpty(&jwtMapClaims.Locale), language.English.String(), jwtMapClaims.Locale),
					Roles:       jwtMapClaims.Roles,
					Permissions: jwtMapClaims.Permissions,
				}
			}

			if r.URL.Query().Has("lang") {
				// Set the AuthUser object in the request context and call the next handler.
				next.ServeHTTP(w, core.WithAuthUser(r, authUser))
			} else {
				localizer := i18n.NewLocalizer(bundle, authUser.GetLocale())
				next.ServeHTTP(w, core.WithLocalizer(core.WithAuthUser(r, authUser), localizer))
			}
		})
	}
}
