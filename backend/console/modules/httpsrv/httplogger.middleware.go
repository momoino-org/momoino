package httpsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5/middleware"
)

// loggerCtxKey is a type alias for a string used as a key for storing a logger in the request context.
type loggerCtxKey string

// HTTPLoggerConfig is a struct used to configure the behavior of the requestLogger middleware.
type HTTPLoggerConfig struct {
	// Tags is a map[string]string used to store additional tags that should be included in the log messages.
	Tags map[string]string

	// IgnoredPaths is a slice of strings representing paths that should be ignored when logging requests.
	// The logger ignores requests with paths that start with any of the ignored paths.
	IgnoredPaths []string
}

// loggerCtxID is a constant string representing the key used to store a logger in the request context.
const loggerCtxID loggerCtxKey = "Logger"

// withLogger is a helper function that adds a logger to the request context.
// It takes an HTTP request and a logger as input and returns a new request with the logger added to its context.
func withLogger(r *http.Request, logger core.Logger) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), loggerCtxID, logger))
}

// getLogger retrieves the logger from the request context.
//
//nolint:ireturn // We don't know what is exactly type of logger, so it is okay to return interface here
func getLogger(r *http.Request) core.Logger {
	logger, _ := r.Context().Value(loggerCtxID).(core.Logger)

	return logger
}

// requestLogger is a middleware function that logs incoming HTTP requests and outgoing HTTP responses.
// It wraps the provided http.Handler and logs the request details, including method, path, headers,
// and elapsed time. It also logs the response details, including status code, bytes written, and headers.
func requestLogger(httpLoggerConfig *HTTPLoggerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.ContainsFunc(httpLoggerConfig.IgnoredPaths, func(ignoredPath string) bool {
				return strings.HasPrefix(r.URL.Path, ignoredPath)
			}) {
				next.ServeHTTP(w, r)
				return
			}

			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			reqCtx := r.Context()

			logger := core.NewStdoutLogger(&core.LoggerConfig{
				RequestHeaderID: "X-Request-Id",
			})

			requestAttrs := []any{
				slog.Attr{Key: "request-id", Value: slog.StringValue(core.GetRequestID(r))},
				slog.Attr{Key: "method", Value: slog.StringValue(r.Method)},
				slog.Attr{Key: "path", Value: slog.StringValue(r.URL.Path)},
				slog.Attr{Key: "headers", Value: slog.AnyValue(r.Header)},
				slog.Attr{Key: "tags", Value: slog.AnyValue(httpLoggerConfig.Tags)},
			}

			logger.Log(
				reqCtx,
				slog.LevelInfo,
				fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path),
				slog.Group("request", requestAttrs...))

			defer func() {
				responseAttrs := []any{
					slog.Attr{Key: "request-id", Value: slog.StringValue(core.GetRequestID(r))},
					slog.Attr{Key: "status", Value: slog.IntValue(ww.Status())},
					slog.Attr{Key: "elapsed", Value: slog.StringValue(time.Since(t1).String())},
					slog.Attr{Key: "bytes", Value: slog.IntValue(ww.BytesWritten())},
					slog.Attr{Key: "headers", Value: slog.AnyValue(ww.Header())},
				}

				logger.Log(
					reqCtx,
					slog.LevelInfo,
					fmt.Sprintf("Response: %s %s", r.Method, r.URL.Path),
					slog.Group("response", responseAttrs...))
			}()

			next.ServeHTTP(ww, withLogger(r, logger))
		})
	}
}
