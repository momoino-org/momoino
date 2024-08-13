package httpsrv

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// withRecoverMiddleware is a middleware function that recovers from panics and logs the panic details.
// It wraps the provided http.Handler and catches any panics that occur during the execution of the handler.
// If a panic occurs, it logs the panic details using the logger retrieved from the request context.
// If the panic is not http.ErrAbortHandler, it writes an HTTP 500 status code to the response writer.
func withRecoverMiddleware(errResponse http.HandlerFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()

				//nolint:errorlint // rvr can be any value, so we cannot use errors.Is
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				if rvr != nil {
					httpLogger, _ := r.Context().Value(loggerCtxID).(*httpLogger)

					httpLogger.logger = httpLogger.logger.With(
						slog.Attr{Key: "panic", Value: slog.StringValue(fmt.Sprintf("%+v", rvr))},
						slog.Attr{Key: "stacktrace", Value: slog.StringValue(string(debug.Stack()))},
					)

					if r.Header.Get("Connection") != "Upgrade" {
						errResponse(w, r)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
