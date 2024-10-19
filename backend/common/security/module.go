package security

import (
	"wano-island/common/core"

	"go.uber.org/fx"
)

func NewSecurityModule() fx.Option {
	return fx.Module("Security module", fx.Provide(
		core.AsRoute(NewGetCsrfTokenHandler),
	))
}
