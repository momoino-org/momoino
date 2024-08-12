package core

import (
	"net/http"

	"go.uber.org/fx"
)

// HTTPRoute is an interface that defines a method for handling HTTP requests.
type HTTPRoute interface {
	http.Handler

	Pattern() string
	IsPrivateRoute() bool
}

func AsRoute(function any) any {
	return fx.Annotate(
		function,
		fx.As(new(HTTPRoute)),
		fx.ResultTags(`group:"http_routes"`),
	)
}
