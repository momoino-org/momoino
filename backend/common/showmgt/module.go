package showmgt

import (
	"wano-island/common/core"

	"go.uber.org/fx"
)

// NewShowMgtModule returns a new Fx module for managing movies.
func NewShowMgtModule() fx.Option {
	return fx.Module(
		"Show management module",
		fx.Provide(
			core.AsRoute(NewCreateMovieHandler),
		),
	)
}
