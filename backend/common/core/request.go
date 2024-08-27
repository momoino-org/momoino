package core

import (
	"net/http"

	"github.com/gorilla/schema"
	"go.uber.org/fx"
)

type requestIDCtxKey string

const RequestIDHeader = "X-Request-Id"
const RequestIDKey requestIDCtxKey = "RequestIDKey"
const AuthorizationHeader = "X-Auth-Access-Token"

func GetRequestID(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}

func NewRequestModule() fx.Option {
	return fx.Module(
		"Request Module",
		fx.Provide(schema.NewDecoder),
	)
}
