package core

import (
	"context"
	"net/http"

	"github.com/gorilla/schema"
	"go.uber.org/fx"
)

type requestIDCtxKey string

const RequestIDKey requestIDCtxKey = "RequestIDKey"

func GetRequestIDInContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}

func GetRequestID(r *http.Request) string {
	return GetRequestIDInContext(r.Context())
}

func NewRequestModule() fx.Option {
	return fx.Module(
		"Request Module",
		fx.Provide(schema.NewDecoder),
	)
}
