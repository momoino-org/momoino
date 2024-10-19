package httpsrv

import (
	"log/slog"
	"net/http"
	"slices"
	"time"
	"wano-island/common/core"
	"wano-island/common/usermgt"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type RouteParams struct {
	fx.In

	Config         core.AppConfig
	Logger         *slog.Logger
	SessionManager *scs.SessionManager
	Routes         []core.HTTPRoute `group:"http_routes"`
	I18nBundle     *i18n.Bundle
	DB             *gorm.DB
	UserRepository usermgt.UserRepository
}

// separatePublicAndPrivateRoutes divides a slice of HTTP routes into two separate slices
// based on their access level, returning one slice for public routes and another for private routes.
//
// Parameters:
//   - routes: A slice of core.HTTPRoute, each containing configuration that determines if it is public or private.
//
// Returns:
//   - []core.HTTPRoute: A slice containing only public routes.
//   - []core.HTTPRoute: A slice containing only private routes.
//
// Each route's access level is determined by inspecting its configuration (`route.Config()`).
// If a route is marked as private, it is added to the private routes slice;
// otherwise, it is included in the public routes slice.
func separatePublicAndPrivateRoutes(routes []core.HTTPRoute) ([]core.HTTPRoute, []core.HTTPRoute) {
	publicRoutes := []core.HTTPRoute{}
	privateRoutes := []core.HTTPRoute{}

	for _, route := range routes {
		routeConfig := route.Config()

		if routeConfig.IsPrivate {
			privateRoutes = append(privateRoutes, route)
		} else {
			publicRoutes = append(publicRoutes, route)
		}
	}

	return publicRoutes, privateRoutes
}

// NewRouter initializes and returns a new HTTP router instance.
func NewRouter(params RouteParams) http.Handler {
	r := chi.NewRouter()

	publicRoutes, privateRoutes := separatePublicAndPrivateRoutes(params.Routes)

	corsConfig := params.Config.GetCorsConfig()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
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
		45: WithCsrfMiddleware(params.Config, params.Logger),
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

		for _, route := range publicRoutes {
			routeConfig := route.Config()

			var wrappedRoute http.Handler = route
			for _, wrapper := range slices.Backward(routeConfig.Wrappers) {
				wrappedRoute = wrapper(wrappedRoute)
			}

			r.Handle(route.Config().Pattern, wrappedRoute)
		}
	})

	// Private routes
	r.Group(func(r chi.Router) {
		middlewares[35] = WithJwtMiddleware(params.I18nBundle, params.Config, params.Logger)
		middlewarePriorities := lo.Keys(middlewares)
		slices.Sort(middlewarePriorities)

		for _, priority := range middlewarePriorities {
			r.Use(middlewares[priority])
		}

		for _, route := range privateRoutes {
			routeConfig := route.Config()

			var wrappedRoute http.Handler = route
			for _, wrapper := range slices.Backward(routeConfig.Wrappers) {
				wrappedRoute = wrapper(wrappedRoute)
			}

			r.Handle(routeConfig.Pattern, wrappedRoute)
		}
	})

	return r
}
