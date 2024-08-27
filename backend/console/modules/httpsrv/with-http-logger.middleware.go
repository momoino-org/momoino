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
	"github.com/samber/lo"
)

// loggerCtxKey is a type alias for a string used as a key for storing a logger in the request context.
type loggerCtxKey int

type customHeader map[string][]string

// HTTPLoggerConfig is a struct used to configure the behavior of the requestLogger middleware.
type HTTPLoggerConfig struct {
	// Tags is a map[string]string used to store additional tags that should be included in the log messages.
	Tags map[string]string

	// IgnoredPaths is a slice of strings representing paths that should be ignored when logging requests.
	// The logger ignores requests with paths that start with any of the ignored paths.
	IgnoredPaths []string
}

// httpLogger is a struct that encapsulates a slog.Logger for logging HTTP requests and responses.
type httpLogger struct {
	logger *slog.Logger
}

// loggerCtxID is a constant string representing the key used to store a logger in the request context.
const loggerCtxID loggerCtxKey = 0

// LogValue is a method that implements the slog.Value interface for the customHeader type.
// It converts the customHeader into a slog.Value, which can be used for structured logging.
// The method redacts sensitive headers (e.g., "Authorization") by replacing their values with "[REDACTED]".
func (c customHeader) LogValue() slog.Value {
	// Define a list of headers to be filtered out and redacted.
	filteredHeaders := []string{
		"Authorization",
		core.AuthorizationHeader,
	}

	// Convert the customHeader into a slice of slog.Attr, filtering out sensitive headers.
	attrs := lo.MapToSlice(c, func(key string, value []string) slog.Attr {
		if slices.Contains(filteredHeaders, key) {
			// If the header is sensitive, redact its value.
			return slog.Any(key, []string{"[REDACTED]"})
		}

		// If the header is not sensitive, include its value as is.
		return slog.Any(key, value)
	})

	// Return the slog.GroupValue containing the filtered and redacted headers.
	return slog.GroupValue(attrs...)
}

// withRequestLoggerMiddleware is a middleware function that logs incoming HTTP requests and outgoing HTTP responses.
// It wraps the provided http.Handler and logs the request details, including method, path, headers,
// and elapsed time. It also logs the response details, including status code, bytes written, and headers.
func withRequestLoggerMiddleware(httpLoggerConfig *HTTPLoggerConfig) func(next http.Handler) http.Handler {
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

			httpLogger := &httpLogger{
				logger: core.NewStdoutLogger(),
			}

			requestAttrs := []any{
				slog.Attr{Key: "request-id", Value: slog.StringValue(core.GetRequestID(r))},
				slog.Attr{Key: "method", Value: slog.StringValue(r.Method)},
				slog.Attr{Key: "path", Value: slog.StringValue(r.URL.Path)},
				slog.Attr{Key: "headers", Value: slog.AnyValue(customHeader(r.Header))},
				slog.Attr{Key: "tags", Value: slog.AnyValue(httpLoggerConfig.Tags)},
			}

			httpLogger.logger.Info(
				fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path),
				slog.Group("request", requestAttrs...))

			defer func() {
				responseAttrs := []any{
					slog.Attr{Key: "request-id", Value: slog.StringValue(core.GetRequestID(r))},
					slog.Attr{Key: "status", Value: slog.IntValue(ww.Status())},
					slog.Attr{Key: "elapsed", Value: slog.StringValue(time.Since(t1).String())},
					slog.Attr{Key: "bytes", Value: slog.IntValue(ww.BytesWritten())},
					slog.Attr{Key: "headers", Value: slog.AnyValue(customHeader(ww.Header()))},
				}

				httpLogger.logger.Info(
					fmt.Sprintf("Response: %s %s", r.Method, r.URL.Path),
					slog.Group("response", responseAttrs...))
			}()

			next.ServeHTTP(ww, r.WithContext(context.WithValue(r.Context(), loggerCtxID, httpLogger)))
		})
	}
}
