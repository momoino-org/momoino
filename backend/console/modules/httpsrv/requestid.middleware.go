package httpsrv

import (
	"context"
	"net/http"
	"wano-island/common/core"

	"github.com/google/uuid"
)

var RequestIDHeader = "X-Request-Id"

func requestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestCtx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
			w.Header().Set(RequestIDHeader, requestID)
		}

		requestCtx = context.WithValue(requestCtx, core.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(requestCtx))
	}

	return http.HandlerFunc(fn)
}
