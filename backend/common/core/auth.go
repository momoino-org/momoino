package core

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type authUserCtx int

type AuthUser struct {
	ID          string
	Locale      string
	Roles       []string
	Permissions []string
}

type JWTCustomClaims struct {
	jwt.RegisteredClaims
	Locale            string   `json:"locale"`
	Roles             []string `json:"roles"`
	Permissions       []string `json:"permissions"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
}

const authUserCtxID authUserCtx = 0

var ErrAuthUserNotFound = errors.New("there is no auth user in the request context")

// WithAuthUser adds the provided AuthUser to the given HTTP request's context.
// This function is useful for passing the authenticated user information throughout the request handling chain.
func WithAuthUser(r *http.Request, user *AuthUser) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authUserCtxID, user))
}

// GetAuthUserFromRequest retrieves the AuthUser from the given HTTP request's context.
// If the AuthUser is found in the context, it will be returned along with a nil error.
// If the AuthUser is not found in the context, it will return nil and an error indicating
// that the AuthUser was not found.
func GetAuthUserFromRequest(r *http.Request) (*AuthUser, error) {
	if authUser, ok := r.Context().Value(authUserCtxID).(*AuthUser); ok {
		return authUser, nil
	}

	return nil, ErrAuthUserNotFound
}

// MustGetAuthUserFromRequest retrieves the AuthUser from the given HTTP request's context.
// If the AuthUser is not found in the context, it will panic with an error message.
func MustGetAuthUserFromRequest(r *http.Request) *AuthUser {
	authUser, err := GetAuthUserFromRequest(r)

	if err != nil {
		panic("Cannot get AuthUser from the request context")
	}

	return authUser
}
