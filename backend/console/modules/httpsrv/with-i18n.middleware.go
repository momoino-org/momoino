package httpsrv

import (
	"net/http"
	"wano-island/common/core"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// withI18nMiddleware is a middleware function that adds internationalization support to an HTTP server.
// It uses the provided i18n.Bundle to create a localizer based on the language specified in the request.
func withI18nMiddleware(bundle *i18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.FormValue("lang")
			acceptLanguages := r.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(bundle, lang, acceptLanguages)

			next.ServeHTTP(w, core.WithLocalizer(r, localizer))
		})
	}
}
