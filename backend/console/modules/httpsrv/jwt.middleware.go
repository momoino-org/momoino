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
)

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

func toAuthUser(mapClaims jwt.MapClaims) (*core.AuthUser, error) {
	authUser := core.AuthUser{}

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

func jwtMiddleware(logger core.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if accessToken == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgNeedToLogin).Build())

				return
			}

			jwtMapClaims := jwt.MapClaims{}
			jwtParser := jwt.NewParser()

			if _, _, err := jwtParser.ParseUnverified(accessToken, &jwtMapClaims); err != nil {
				logger.ErrorContext(r.Context(), "Cannot parse jwt", slog.Any("details", err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

				return
			}

			authUser, err := toAuthUser(jwtMapClaims)
			if err != nil {
				logger.ErrorContext(r.Context(), "Cannot convert jwt claims to AuthUser", slog.Any("details", err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgCannotProcessYourLogin).Build())

				return
			}

			next.ServeHTTP(w, core.WithAuthUser(r, authUser))
		})
	}
}
