package httpsrv

import "go.uber.org/fx"

func NewHTTPServerModule() fx.Option {
	return fx.Module(
		"HTTP Server Module",
		fx.Provide(
			NewRouter,
			newHTTPServer,
		),
	)
}
