package httpsrv

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/gorilla/csrf"
)

// WithCsrfMiddleware returns an HTTP middleware that provides CSRF (Cross-Site Request Forgery)
// protection for incoming requests. It utilizes a CSRF token mechanism to ensure the request
// originates from a trusted source, and configures several security-related settings to align with
// application requirements.
//
// The middleware is configured based on the provided `core.AppConfig`, which supplies the
// necessary CSRF secret key and HTTPS status. If the CSRF validation fails, an error handler logs
// the reason and responds with a 403 Forbidden status and a standardized error message.
//
// Parameters:
//   - config: core.AppConfig that contains application-specific settings, including the CSRF
//     secret key and HTTPS enforcement.
//   - logger: slog.Logger for logging CSRF failure reasons.
//
// Returns:
//   - A middleware function that wraps an HTTP handler, adding CSRF protection.
func WithCsrfMiddleware(
	config core.AppConfig,
	logger *slog.Logger,
) func(next http.Handler) http.Handler {
	return csrf.Protect(
		config.GetSecretKey(),
		csrf.TrustedOrigins([]string{"localhost"}),
		csrf.Secure(config.IsHTTPS()),
		csrf.CookieName(core.CsrfCookie),
		csrf.Path("/"),
		csrf.MaxAge(0),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.ErrorContext(r.Context(), csrf.FailureReason(r).Error())
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInvalidCsrf).Build())
		})),
	)
}
