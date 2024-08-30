package httpsrv

import (
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
)

type RouteParams struct {
	fx.In

	Config     core.AppConfig
	Logger     *slog.Logger
	Routes     []core.HTTPRoute `group:"http_routes"`
	I18nBundle *i18n.Bundle
}

// NewRouter initializes and returns a new HTTP router instance.
func NewRouter(params RouteParams) http.Handler {
	r := chi.NewRouter()

	publicRoutes := []core.HTTPRoute{}
	privateRoutes := []core.HTTPRoute{}

	for _, route := range params.Routes {
		if route.IsPrivateRoute() {
			privateRoutes = append(privateRoutes, route)
		} else {
			publicRoutes = append(publicRoutes, route)
		}
	}

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.CleanPath,
			withRequestIDMiddleware,
			withI18nMiddleware(params.I18nBundle),
			withRequestLoggerMiddleware(params.Config, &HTTPLoggerConfig{
				Silent: params.Config.IsTesting(),
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
			withRecoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())
			}),
			middleware.Timeout(time.Minute),
		)

		for _, route := range publicRoutes {
			r.Handle(route.Pattern(), route)
		}
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.CleanPath,
			withRequestIDMiddleware,
			withI18nMiddleware(params.I18nBundle),
			withJwtMiddleware(params.I18nBundle, params.Logger),
			withRequestLoggerMiddleware(params.Config, &HTTPLoggerConfig{
				Silent: params.Config.IsTesting(),
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
			withRecoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())
			}),
			middleware.Timeout(time.Minute),
		)

		for _, route := range privateRoutes {
			r.Handle(route.Pattern(), route)
		}
	})

	return r
}
