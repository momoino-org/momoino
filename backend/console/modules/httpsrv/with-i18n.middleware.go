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
			authUser := core.GetAuthUserFromRequest(r)
			lang := r.URL.Query().Get("lang")

			if lang != "" && authUser.Locale != lang {
				authUser.Locale = lang
			}

			localizer := i18n.NewLocalizer(bundle, lang, authUser.Locale)

			next.ServeHTTP(w, core.WithLocalizer(r, localizer))
		})
	}
}
