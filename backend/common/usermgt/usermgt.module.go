package usermgt

import (
	"wano-island/common/core"

	"go.uber.org/fx"
)

// NewUserMgtModule is a function that returns an fx.Option for the User management module.
// This module provides services and handlers related to user management, including login, profile,
// password change, and OAuth2 authentication.
func NewUserMgtModule() fx.Option {
	return fx.Module(
		"User management module",
		fx.Provide(
			fx.Annotate(NewUserService, fx.As(new(UserService))),
			fx.Annotate(NewUserRepository, fx.As(new(UserRepository))),
			core.AsRoute(NewLoginHandler),
			core.AsRoute(NewProfileHandler),
			core.AsRoute(NewChangePasswordHandler),
			core.AsRoute(NewRenewSessionHandler),
			core.AsRoute(NewCreateLoginSessionHandler),

			// OAuth2
			fx.Annotate(NewGoogleProvider, fx.As(new(OAuth2Provider)), fx.ResultTags(`name:"google_provider"`)),
			fx.Annotate(NewOAuth2ProviderRepository, fx.As(new(OAuth2ProviderRepository))),
			core.AsRoute(NewCreateOAuth2Provider),
			core.AsRoute(NewOAuth2LoginHandler),
			core.AsRoute(NewOAuth2LoginCallbackHandler),
			core.AsRoute(NewGetOAuth2ProvidersHandler),
		),
	)
}
