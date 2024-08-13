package httpsrv

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/google/uuid"
)

// withRequestIDMiddleware is a middleware function for HTTP requests that adds a unique request ID to each request.
// If a request already contains a request ID in its header, it will be used. Otherwise, a new UUID will be generated.
// The request ID will be added to the response header and the request context.
// The middleware function wraps the provided http.Handler and adds the request ID functionality.
func withRequestIDMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestCtx := r.Context()
		requestID := r.Header.Get(core.RequestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set(core.RequestIDHeader, requestID)
			w.Header().Set(core.RequestIDHeader, requestID)
		}

		requestCtx = core.WithLogAttr(requestCtx, slog.String("request-id", requestID))
		next.ServeHTTP(w, r.WithContext(requestCtx))
	}

	return http.HandlerFunc(fn)
}
