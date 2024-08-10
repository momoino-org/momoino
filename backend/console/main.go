package main

import (
	"embed"
	"net/http"
	"time"
	"wano-island/common/core"
	"wano-island/common/usermgt"
	"wano-island/console/modules/filesystem"
	"wano-island/console/modules/httpsrv"
	"wano-island/console/modules/swagger"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

//go:embed static
var staticFiles embed.FS

func main() {
	time.Local = time.UTC

	app := fx.New(
		// Common
		core.NewConfigModule(),
		core.NewLoggerModuleWithConfig(&core.LoggerConfig{
			RequestHeaderID: middleware.RequestIDHeader,
		}),
		core.NewRequestModule(),
		core.NewDatabaseModule(),
		usermgt.NewUserMgtModule(),

		// Console
		filesystem.NewFileSystemModule(staticFiles),
		swagger.NewSwaggerModule(),
		httpsrv.NewHTTPServerModule(),

		// Start web server
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
