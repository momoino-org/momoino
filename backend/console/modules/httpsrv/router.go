package httpsrv

import (
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/samber/lo"
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

	httplogger := httplog.NewLogger("x-operation/console", httplog.Options{
		JSON:            params.Config.IsDevelopment(),
		LogLevel:        lo.Ternary(params.Config.IsProduction(), slog.LevelInfo, slog.LevelDebug),
		RequestHeaders:  true,
		TimeFieldFormat: time.RFC3339Nano,
		TimeFieldName:   "time",
		Tags: map[string]string{
			"version": params.Config.GetAppVersion(),
			"mode":    params.Config.GetMode(),
			"rev":     params.Config.GetRevision(),
		},
	})

	r.Use(
		middleware.CleanPath,
		middleware.RequestID,
		httplog.Handler(httplogger),
		middleware.Recoverer,
		middleware.Timeout(time.Minute),
	)

	for _, route := range params.Routes {
		r.Handle(route.Pattern(), route)
	}

	return r
}
