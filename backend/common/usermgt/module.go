package usermgt

import (
	"wano-island/common/core"

	"go.uber.org/fx"
)

func NewUserMgtModule() fx.Option {
	return fx.Module(
		"User management module",
		fx.Provide(
			fx.Annotate(NewUserService, fx.As(new(UserService))),
			fx.Annotate(NewUserRepository, fx.As(new(UserRepository))),
			core.AsRoute(NewLoginHandler),
			core.AsRoute(NewProfileHandler),
			core.AsRoute(NewChangePasswordHandler),
		),
	)
}
