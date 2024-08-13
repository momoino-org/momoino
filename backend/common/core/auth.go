package core

import (
	"context"
	"net/http"
)

type authUserCtx int

type AuthUser struct {
	ID          string
	Roles       []string
	Permissions []string
}

const authUserCtxID authUserCtx = 0

func WithAuthUser(r *http.Request, user *AuthUser) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authUserCtxID, user))
}

func GetAuthUserFromRequest(r *http.Request) (*AuthUser, bool) {
	authUser, ok := r.Context().Value(authUserCtxID).(*AuthUser)

	return authUser, ok
}
