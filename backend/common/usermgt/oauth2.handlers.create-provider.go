package usermgt

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

// createOAuth2ProviderHandler handles the creation of new OAuth2 providers.
type createOAuth2ProviderHandler struct {
	logger                   *slog.Logger
	validator                *validator.Validate
	universalTranslator      *ut.UniversalTranslator
	oauth2ProviderRepository OAuth2ProviderRepository
}

// CreateOAuth2ProviderHandlerParams defines the dependencies required to initialize a createOAuth2ProviderHandler.
type CreateOAuth2ProviderHandlerParams struct {
	fx.In
	Logger                   *slog.Logger
	Validator                *validator.Validate
	UniversalTranslator      *ut.UniversalTranslator
	OAuth2ProviderRepository OAuth2ProviderRepository
}

// CreateOAuth2ProviderRequestBody holds the request body for creating a new OAuth2 provider.
// It includes fields that must be validated to ensure the integrity of the data being saved.
type CreateOAuth2ProviderRequestBody struct {
	// The name of the OAuth2 provider.
	Provider string `json:"provider" validate:"required"`

	// The client ID for the OAuth2 provider.
	ClientID string `json:"clientID" validate:"required"`

	// The client secret for the OAuth2 provider.
	ClientSecret string `json:"clientSecret" validate:"required"`

	// The redirect URL for OAuth2 authentication.
	RedirectURL string `json:"redirectUrl" validate:"required"`

	// The scopes that the provider will request.
	Scopes []string `json:"scopes" validate:"required"`

	// A flag indicating whether the provider is enabled.
	IsEnabled *bool `json:"isEnabled" validate:"required"`
}

// Ensure that createOAuth2ProviderHandler implements the core.HTTPRoute interface.
// This guarantees that createOAuth2ProviderHandler conforms to the expected route handler structure.
var _ core.HTTPRoute = (*createOAuth2ProviderHandler)(nil)

// NewCreateOAuth2Provider initializes and returns a new instance of createOAuth2ProviderHandler.
// It sets up the handler with the required dependencies.
//
// Params:
//   - p: The parameters required for initializing the handler, including logger, validator, translator, and repository.
//
// Returns:
//   - A new createOAuth2ProviderHandler instance.
func NewCreateOAuth2Provider(p CreateOAuth2ProviderHandlerParams) *createOAuth2ProviderHandler {
	return &createOAuth2ProviderHandler{
		logger:                   p.Logger,
		validator:                p.Validator,
		universalTranslator:      p.UniversalTranslator,
		oauth2ProviderRepository: p.OAuth2ProviderRepository,
	}
}

func (h *createOAuth2ProviderHandler) Pattern() string {
	return "POST /api/v1/providers"
}

func (h *createOAuth2ProviderHandler) IsPrivateRoute() bool {
	return true
}

// ServeHTTP processes incoming HTTP requests for creating new OAuth2 providers.
// It decodes the request body, validates the data, and invokes the repository to create the provider.
// If successful, it returns a 201 Created response with the new provider details; otherwise, it handles
// errors accordingly.
// Params:
//   - w: The HTTP response writer to send the response.
//   - r: The HTTP request containing the data for the new OAuth2 provider.
//
// Behavior:
//   - If the request body cannot be decoded, it returns a 400 Bad Request error.
//   - If validation fails, it returns a 400 Bad Request error with validation messages.
//   - If the creation of the provider fails, it logs the error and returns a 500 Internal Server Error.
//   - On successful creation, it returns a 201 Created response with the created provider's data.
func (h *createOAuth2ProviderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseBuilder := core.NewResponseBuilder(r)

	// Decode the request body
	var requestBody CreateOAuth2ProviderRequestBody
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	// Validate the request body
	if err := h.validator.Struct(requestBody); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.
			MessageID(core.MsgValidationFailed).
			Data(core.TranslateValidationErrors(r, h.universalTranslator, err)).
			Build())

		return
	}

	// Create the OAuth2 provider
	oauth2Provider, err := h.oauth2ProviderRepository.Create(
		r.Context(),
		CreateOAuth2ProviderParams{
			Provider:     requestBody.Provider,
			ClientID:     requestBody.ClientID,
			ClientSecret: requestBody.ClientSecret,
			RedirectURL:  requestBody.RedirectURL,
			Scopes:       requestBody.Scopes,
			IsEnabled:    *requestBody.IsEnabled,
			CreatedBy:    core.MustGetAuthUserFromRequest(r),
		},
	)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot create the oauth2 provider", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Return the created provider
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, responseBuilder.MessageID(core.MsgSuccess).Data(OAuth2ProviderDTO{
		ID:        oauth2Provider.ClientID,
		Provider:  oauth2Provider.Provider,
		IsEnabled: oauth2Provider.IsEnabled,
		CreatedAt: oauth2Provider.CreatedAt,
		CreatedBy: oauth2Provider.CreatedBy,
	}).Build())
}
