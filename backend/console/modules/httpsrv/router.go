package httpsrv

import (
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

type routeParams struct {
	fx.In

	Config core.AppConfig
	Logger core.Logger
	Routes []core.HTTPRoute `group:"http_routes"`
}

// newRouter initializes and returns a new HTTP router instance.
func newRouter(params routeParams) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.CleanPath,
		requestID,
		requestLogger(&HTTPLoggerConfig{
			IgnoredPaths: []string{
				"/swagger",
				"/static",
			},
			Tags: map[string]string{
				"version": params.Config.GetAppVersion(),
				"mode":    params.Config.GetMode(),
				"rev":     params.Config.GetRevision(),
			},
		}),
		//nolint:mnd // I don't think we need to named this number here
		middleware.Compress(5),
		httpRecover,
		middleware.Timeout(time.Minute),
	)

	for _, route := range params.Routes {
		r.Handle(route.Pattern(), route)
	}

	return r
}
