package httpsrv

import (
	"log/slog"
	"net/http"
	"strings"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

// extractAccessTokenFromRequest extracts the access token from the HTTP request header or cookie.
func extractAccessTokenFromRequest(req *http.Request) string {
	// Extract the access token from the request header.
	accessToken := strings.TrimPrefix(req.Header.Get(core.AuthorizationHeader), "Bearer ")

	// Extract the access token from the cookie if it doesn not exist in the request header.
	if len(accessToken) == 0 {
		if cookie, err := req.Cookie(core.IdentityCookie); err == nil {
			accessToken = cookie.Value
		}
	}

	return strings.TrimSpace(accessToken)
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

			if _, err := jwtParser.ParseWithClaims(accessToken, &jwtMapClaims, func(t *jwt.Token) (interface{}, error) {
				return config.GetJWTConfig().PublicKey, nil
			}); err != nil {
				logger.ErrorContext(r.Context(), "Cannot parse jwt", core.DetailsLogAttr(err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

				return
			}

			// Convert the JWT claims into an AuthUser object.
			authUser := &core.AuthenticatedUser{
				ID:          jwtMapClaims.Subject,
				Username:    jwtMapClaims.PreferredUsername,
				Email:       jwtMapClaims.Email,
				FamilyName:  jwtMapClaims.FamilyName,
				GivenName:   jwtMapClaims.GivenName,
				Locale:      lo.Ternary(lo.IsEmpty(&jwtMapClaims.Locale), language.English.String(), jwtMapClaims.Locale),
				Roles:       jwtMapClaims.Roles,
				Permissions: jwtMapClaims.Permissions,
			}

			if r.URL.Query().Has("lang") {
				// Set the AuthUser object in the request context and call the next handler.
				next.ServeHTTP(w, core.WithAuthUser(r, authUser))
			} else {
				localizer := i18n.NewLocalizer(bundle, authUser.Locale)
				next.ServeHTTP(w, core.WithLocalizer(core.WithAuthUser(r, authUser), localizer))
			}
		})
	}
}
