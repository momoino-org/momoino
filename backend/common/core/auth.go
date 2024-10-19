package core

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type authUserCtx int

type PrincipalUser interface {
	GetID() uuid.UUID
	GetUsername() string
	GetEmail() string
	GetGivenName() string
	GetFamilyName() string
	GetLocale() string
	GetRoles() []string
	GetPermissions() []string
}

type AuthenticatedUser struct {
	ID          string
	Username    string
	Email       string
	GivenName   string
	FamilyName  string
	Locale      string
	Roles       []string
	Permissions []string
}

type JWTCustomClaims struct {
	jwt.RegisteredClaims
	SessionID         string   `json:"sid"`
	Locale            string   `json:"locale"`
	Roles             []string `json:"roles"`
	Permissions       []string `json:"permissions"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
}

const authUserCtxID authUserCtx = 0
const IdentityCookie = "MOMOINO_IDENTITY"
const SessionCookie = "MOMOINO_SESSION"
const LoginSessionCookie = "MOMOINO_LOGIN_SESSION"
const CsrfCookie = "MOMOINO_CSRF"

var _ PrincipalUser = (*AuthenticatedUser)(nil)

var ErrAuthUserNotFound = errors.New("there is no auth user in the request context")

// WithAuthUser adds the provided AuthUser to the given HTTP request's context.
// This function is useful for passing the authenticated user information throughout the request handling chain.
func WithAuthUser(r *http.Request, user PrincipalUser) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authUserCtxID, user))
}

// GetAuthUserFromRequest retrieves the AuthUser from the given HTTP request's context.
// If the AuthUser is found in the context, it will be returned along with a nil error.
// If the AuthUser is not found in the context, it will return nil and an error indicating
// that the AuthUser was not found.
//
//nolint:ireturn // I need to return interface at here
func GetAuthUserFromRequest(r *http.Request) (PrincipalUser, error) {
	if authUser, ok := r.Context().Value(authUserCtxID).(PrincipalUser); ok {
		return authUser, nil
	}

	return nil, ErrAuthUserNotFound
}

// MustGetAuthUserFromRequest retrieves the AuthUser from the given HTTP request's context.
// If the AuthUser is not found in the context, it will panic with an error message.
//
//nolint:ireturn // I need to return interface at here
func MustGetAuthUserFromRequest(r *http.Request) PrincipalUser {
	authUser, err := GetAuthUserFromRequest(r)

	if err != nil {
		panic("Cannot get AuthUser from the request context")
	}

	return authUser
}

func (u AuthenticatedUser) GetID() uuid.UUID {
	return uuid.MustParse(u.ID)
}

func (u AuthenticatedUser) GetUsername() string {
	return u.Username
}

func (u AuthenticatedUser) GetEmail() string {
	return u.Email
}

func (u AuthenticatedUser) GetGivenName() string {
	return u.GivenName
}

func (u AuthenticatedUser) GetFamilyName() string {
	return u.FamilyName
}

func (u AuthenticatedUser) GetLocale() string {
	if len(u.Locale) == 0 {
		return language.English.String()
	}

	return u.Locale
}

func (u AuthenticatedUser) GetRoles() []string {
	return []string{}
}

func (u AuthenticatedUser) GetPermissions() []string {
	return []string{}
}
