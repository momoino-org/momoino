package httpsrv

import (
	"context"
	"net/http"
	"wano-island/common/core"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func i18nMiddleware(bundle *i18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.FormValue("lang")
			acceptLanguages := r.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(bundle, lang, acceptLanguages)

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), core.LocalizerCtxID, localizer)))
		})
	}
}
