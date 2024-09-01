package httpsrv

import (
	"log/slog"
	"net/http"
	"slices"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Use this to allow specific origin hosts
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	middlewares := map[int]func(http.Handler) http.Handler{
		0:  middleware.CleanPath,
		10: withRequestIDMiddleware,
		20: withRequestLoggerMiddleware(&HTTPLoggerConfig{
			CustomLogger: func() *slog.Logger {
				if params.Config.IsTesting() {
					return core.NewNoopLogger()
				}

				return core.NewStdoutLogger(params.Config)
			},
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
		30: WithI18nMiddleware(params.I18nBundle),
		//nolint:mnd // I don't think we need to named this number here
		40: middleware.Compress(5),
		50: withRecoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())
		}),
		60: middleware.Timeout(time.Minute),
	}

	// Public routes
	r.Group(func(r chi.Router) {
		middlewarePriorities := lo.Keys(middlewares)
		slices.Sort(middlewarePriorities)

		for _, priority := range middlewarePriorities {
			r.Use(middlewares[priority])
		}

		r.Use(lo.Values(middlewares)...)

		for _, route := range publicRoutes {
			r.Handle(route.Pattern(), route)
		}
	})

	// Private routes
	r.Group(func(r chi.Router) {
		middlewares[35] = WithJwtMiddleware(params.I18nBundle, params.Logger)
		middlewarePriorities := lo.Keys(middlewares)
		slices.Sort(middlewarePriorities)

		for _, priority := range middlewarePriorities {
			r.Use(middlewares[priority])
		}

		for _, route := range privateRoutes {
			r.Handle(route.Pattern(), route)
		}
	})

	return r
}
