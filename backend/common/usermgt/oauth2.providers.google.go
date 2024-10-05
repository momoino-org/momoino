package usermgt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2API "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// googleProvider is the implementation of the OAuth2Provider interface for Google's OAuth2 service.
// It handles the specific operations needed to authenticate and interact with Google's OAuth2 system.
type googleProvider struct {
	logger                   *slog.Logger
	config                   core.AppConfig
	db                       *gorm.DB
	schemaDecoder            *schema.Decoder
	userService              UserService
	userRepository           UserRepository
	oauth2ProviderRepository OAuth2ProviderRepository
}

// GoogleProviderHandlerParams defines the dependencies required to initialize a googleProvider.
type GoogleProviderHandlerParams struct {
	fx.In
	Logger                   *slog.Logger
	Config                   core.AppConfig
	DB                       *gorm.DB
	SchemaDecoder            *schema.Decoder
	UserService              UserService
	UserRepository           UserRepository
	Oauth2ProviderRepository OAuth2ProviderRepository
}

// AuthorizeRequestQueryParams defines the required query parameters for initiating the OAuth2 authorization request.
type AuthorizeRequestQueryParams struct {
	// A unique string to maintain state between the request and callback, and prevent CSRF attacks.
	State string `schema:"state,required"`

	// A challenge for securing the authorization code exchange using PKCE (Proof Key for Code Exchange).
	CodeChallenge string `schema:"codeChallenge,required"`
}

// CallbackRequestQueryParams defines the required query parameters for the OAuth2 callback endpoint.
type CallbackRequestQueryParams struct {
	// The state parameter returned by Google to maintain state between requests.
	State string `schema:"state,required"`

	// The authorization code issued by Google, to be exchanged for an access token.
	Code string `schema:"code,required"`

	// The verifier corresponding to the code challenge, for securing the token exchange with PKCE.
	Verifier string `schema:"verifier,required"`
}

// Ensure that googleProvider implements the OAuth2Provider interface.
// This line will cause a compile-time error if googleProvider doesn't satisfy the interface.
var _ OAuth2Provider = (*googleProvider)(nil)

// NewGoogleProvider creates and returns a new googleProvider instance.
// It initializes the provider with the necessary dependencies for handling OAuth2 interactions with Google.
//
// Params:
//   - p: Dependencies for initializing the Google provider, injected via the Fx framework.
//
// Returns:
//   - A new googleProvider instance that can handle OAuth2 authentication flows with Google.
func NewGoogleProvider(p GoogleProviderHandlerParams) *googleProvider {
	return &googleProvider{
		logger:                   p.Logger,
		config:                   p.Config,
		db:                       p.DB,
		schemaDecoder:            p.SchemaDecoder,
		userService:              p.UserService,
		userRepository:           p.UserRepository,
		oauth2ProviderRepository: p.Oauth2ProviderRepository,
	}
}

// getUserInfo retrieves user information from the Google OAuth2 service.
//
// Parameters:
//   - ctx: The context for the request.
//   - httpClient: The HTTP client to be used for making the request.
//
// Returns:
//   - userInfo: The user information retrieved from the Google API.
//   - err: An error if any occurred during the request or parsing the response.
func (p *googleProvider) getUserInfo(ctx context.Context, httpClient *http.Client) (*oauth2API.Userinfo, error) {
	oauth2Service, err := oauth2API.NewService(ctx, option.WithTelemetryDisabled(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("cannot create OAuth2Service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("couldn't get user info: %w", err)
	}

	return userInfo, nil
}

// getProvider retrieves the OAuth2ProviderModel for Google from the database.
// It checks if the provider is enabled and returns the model if available.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - provider: The OAuth2ProviderModel for Google if found and enabled.
//   - err: An error if any occurred during the database query or if the provider is not enabled.
func (p *googleProvider) getProvider(ctx context.Context) (*OAuth2ProviderModel, error) {
	provider, err := p.oauth2ProviderRepository.Get(ctx, GoogleProviderName)

	if err != nil {
		return nil, err
	}

	if !provider.IsEnabled {
		return nil, fmt.Errorf("OAuth2 provider (%v) is not enabled", provider.Provider)
	}

	return provider, nil
}

// buildOAuth2Config creates and returns an OAuth2 configuration for Google.
// The configuration is based on the provided OAuth2ProviderModel.
//
// Parameters:
//   - provider: The OAuth2ProviderModel containing the necessary client ID, client secret, redirect URL, and scopes.
//
// Returns:
//   - An OAuth2 configuration for Google, which can be used to initiate the OAuth2 authorization flow.
func (p *googleProvider) buildOAuth2Config(provider *OAuth2ProviderModel) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
		Scopes:       provider.Scopes,
		Endpoint:     google.Endpoint,
	}
}

// AuthorizeHandler handles the OAuth2 authorization flow for Google.
// It retrieves the OAuth2 provider from the database, builds the OAuth2 configuration,
// decodes the query parameters, generates the authorization URL, and redirects the user to Google.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//
// Returns:
//   - No return value. The function writes the HTTP response directly.
func (p *googleProvider) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	providerModel, err := p.getProvider(reqCtx)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.ErrorContext(reqCtx,
				"Something went wrong when getting the OAuth2 provider from the database.",
				core.DetailsLogAttr(err))
		}

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgOAuth2ProviderUnavailable).Build())

		return
	}

	oauth2Config := p.buildOAuth2Config(providerModel)

	var params AuthorizeRequestQueryParams
	if err := p.schemaDecoder.Decode(&params, r.URL.Query()); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		params.State,
		oauth2.AccessTypeOnline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge))

	http.Redirect(w, r, authCodeURL, http.StatusTemporaryRedirect)
}

// CallbackHandler handles the OAuth2 callback flow for Google.
// It retrieves the OAuth2 provider from the database, builds the OAuth2 configuration,
// decodes the query parameters, exchanges the authorization code for an access token,
// retrieves user information from the Google API, creates or updates the user in the database,
// generates a JWT for the user, and sets the JWT as authentication cookies in the HTTP response.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
//
// Return Value:
//   - No return value. The function writes the HTTP response directly.
func (p *googleProvider) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	providerModel, err := p.getProvider(r.Context())
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.ErrorContext(reqCtx,
				"Something went wrong when getting the OAuth2 provider from the database.",
				core.DetailsLogAttr(err))
		}

		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgOAuth2ProviderUnavailable).Build())

		return
	}

	oauth2Config := p.buildOAuth2Config(providerModel)

	var params CallbackRequestQueryParams
	if decodeErr := p.schemaDecoder.Decode(&params, r.URL.Query()); decodeErr != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	oauth2Token, exchangeErr := oauth2Config.Exchange(reqCtx, params.Code, oauth2.VerifierOption(params.Verifier))
	if exchangeErr != nil {
		p.logger.ErrorContext(reqCtx, "Cannot convert an authorization code into a token.", core.DetailsLogAttr(exchangeErr))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	client := oauth2Config.Client(reqCtx, oauth2Token)

	googleUser, getUserInfoErr := p.getUserInfo(reqCtx, client)
	if getUserInfoErr != nil {
		p.logger.ErrorContext(reqCtx,
			"Cannot get user info from the google provider.",
			core.DetailsLogAttr(getUserInfoErr))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	userModel, err := p.userRepository.FirstOrCreateUser(reqCtx, p.db, CreateUserParams{
		Username:      googleUser.Email,
		Email:         googleUser.Email,
		VerifiedEmail: lo.Ternary(googleUser.VerifiedEmail == nil, false, *googleUser.VerifiedEmail),
		FirstName:     googleUser.GivenName,
		LastName:      googleUser.FamilyName,
		CreatedBy:     core.NewSystemUser(),
		LinkedProviders: []LinkedProvider{
			{
				Provider: *providerModel,
				OpenID:   googleUser.Id,
			},
		},
	})
	if err != nil {
		p.logger.ErrorContext(reqCtx, "Something went wrong when getting/creating the oauth2 user.", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	jwt, err := p.userService.GenerateJWT(*userModel)
	if err != nil {
		p.logger.ErrorContext(r.Context(), "Something went wrong when creating jwt.", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	p.userService.SetAuthCookies(w, *jwt)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Build())
}
