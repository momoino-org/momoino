package core

import (
	"net/http"

	"go.uber.org/fx"
)

// HTTPRoute is an interface that defines a method for handling HTTP requests.
type HTTPRoute interface {
	http.Handler

	Config() *HTTPRouteConfig
}

type HTTPRouteConfig struct {
	Pattern   string
	IsPrivate bool
	Wrappers  []func(http.Handler) http.Handler
}

func AsRoute(function any) any {
	return fx.Annotate(
		function,
		fx.As(new(HTTPRoute)),
		fx.ResultTags(`group:"http_routes"`),
	)
}
