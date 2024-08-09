package main

import (
	"net/http"
	"wano-island/common/core"
	"wano-island/console/modules/filesystem"
	"wano-island/console/modules/httpsrv"
	"wano-island/console/modules/swagger"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Common
		core.NewConfigModule(),
		core.NewLoggerModuleWithConfig(&core.LoggerConfig{
			RequestHeaderID: middleware.RequestIDHeader,
		}),

		// Console
		filesystem.NewFileSystemModule(),
		swagger.NewSwaggerModule(),
		httpsrv.NewHTTPServerModule(),

		// Start web server
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
