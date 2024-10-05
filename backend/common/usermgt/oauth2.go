package usermgt

import (
	"net/http"
)

// OAuth2Provider is an interface that defines methods for handling OAuth2 authentication.
type OAuth2Provider interface {
	// AuthorizeHandler redirects the user to the OAuth2 provider's authorization page.
	AuthorizeHandler(w http.ResponseWriter, r *http.Request)

	// CallbackHandler handles the OAuth2 provider's callback after user authorization.
	CallbackHandler(w http.ResponseWriter, r *http.Request)
}

const GoogleProviderName = "google"
