package httpsrv

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
	"wano-island/common/core"

	"go.uber.org/fx"
)

// newHTTPServer initializes and returns a new HTTP server instance.
func newHTTPServer(
	appLifeCycle fx.Lifecycle,
	httpHandler http.Handler,
	logger core.Logger,
) *http.Server {
	const readHeaderTimeout = 5 * time.Second

	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           httpHandler,
	}

	appLifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func(ctx context.Context) {
				err := srv.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					logger.ErrorContext(ctx, "Cannot start http server", slog.Any("details", err))
					os.Exit(1)
				}
			}(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
