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

	"go.uber.org/fx"
)

//go:embed static
var staticFiles embed.FS

//go:embed resources
var resourceFS embed.FS

func main() {
	time.Local = time.UTC

	app := fx.New(
		// Common
		core.NewI18nModule(resourceFS),
		core.NewValidationModule(),
		core.NewConfigModule(),
		core.NewLoggerModuleWithConfig(),
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
