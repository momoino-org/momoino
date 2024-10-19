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
			core.AsRoute(NewProfileHandler),
		),
	)
}
