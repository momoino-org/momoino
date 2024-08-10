package httpsrv

import (
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
)

type routeParams struct {
	fx.In

	Config     core.AppConfig
	Logger     core.Logger
	Routes     []core.HTTPRoute `group:"http_routes"`
	I18nBundle *i18n.Bundle
}

// newRouter initializes and returns a new HTTP router instance.
func newRouter(params routeParams) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.CleanPath,
		requestID,
		i18nMiddleware(params.I18nBundle),
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
		httpRecover(func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID("U-0001").Build())
		}),
		middleware.Timeout(time.Minute),
	)

	for _, route := range params.Routes {
		r.Handle(route.Pattern(), route)
	}

	return r
}
