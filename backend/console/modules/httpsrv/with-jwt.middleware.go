package httpsrv

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

// convertInterfaceSliceToStringSlice converts a slice of interface{} to a slice of string.
// If any element is not a string, it returns an error indicating the index of the invalid element.
func convertInterfaceSliceToStringSlice(data []interface{}) ([]string, error) {
	result := make([]string, len(data))

	for i, v := range data {
		str, ok := v.(string)

		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}

		result[i] = str
	}

	return result, nil
}

// toAuthUser converts JWT claims into an AuthUser object.
// It extracts the user ID, roles, and permissions from the claims and constructs an AuthUser object.
// If any required claim is missing or cannot be converted to the expected type, an error is returned.
func toAuthUser(mapClaims jwt.MapClaims) (*core.AuthUser, error) {
	authUser := core.AuthUser{
		Locale: language.English.String(),
	}

	userID, err := mapClaims.GetSubject()
	if err != nil {
		return nil, errors.New("cannot get subject from claims")
	} else {
		authUser.ID = userID
	}

	if !lo.HasKey(mapClaims, "roles") {
		return nil, errors.New("cannot get roles from claims")
	}

	if !lo.HasKey(mapClaims, "permissions") {
		return nil, errors.New("cannot get permissions from claims")
	}

	roles, err := convertInterfaceSliceToStringSlice(mapClaims["roles"].([]interface{}))
	if err == nil {
		authUser.Roles = roles
	} else {
		return nil, err
	}

	permissions, err := convertInterfaceSliceToStringSlice(mapClaims["permissions"].([]interface{}))
	if err == nil {
		authUser.Permissions = permissions
	} else {
		return nil, err
	}

	return &authUser, nil
}

// withJwtMiddleware is a middleware function that handles JWT authentication for HTTP requests.
// It extracts the access token from the request header, verifies it, and converts the claims into an AuthUser object.
// If the access token is missing or invalid, it returns an appropriate HTTP response.
// Otherwise, it sets the AuthUser object in the request context and calls the next handler.
func withJwtMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the access token from the request header.
			accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			// If the access token is missing, return an unauthorized response.
			if accessToken == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgNeedToLogin).Build())

				return
			}

			// Parse the access token without verifying it to extract the JWT claims.
			jwtMapClaims := jwt.MapClaims{}
			jwtParser := jwt.NewParser()

			if _, _, err := jwtParser.ParseUnverified(accessToken, &jwtMapClaims); err != nil {
				logger.ErrorContext(r.Context(), "Cannot parse jwt", slog.Any("details", err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

				return
			}

			// Convert the JWT claims into an AuthUser object.
			authUser, err := toAuthUser(jwtMapClaims)
			if err != nil {
				logger.ErrorContext(r.Context(), "Cannot convert jwt claims to AuthUser", slog.Any("details", err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

				return
			}

			// Set the AuthUser object in the request context and call the next handler.
			next.ServeHTTP(w, core.WithAuthUser(r, authUser))
		})
	}
}
