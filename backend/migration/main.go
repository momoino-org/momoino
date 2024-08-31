package main

import (
	"context"
	"log/slog"
	"wano-island/common/core"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions"

	"go.uber.org/fx"
)

// main is the entry point of the migration application.
// It initializes the necessary dependencies and runs the database migration process.
func main() {
	fx.New(
		core.NewConfigModule(),
		core.NewLoggerModuleWithConfig(),
		core.NewDatabaseModule(),
		versions.NewDBMigrationModule(),
		fx.Invoke(func(
			appLifeCycle fx.Lifecycle,
			logger *slog.Logger,
			shutdowner fx.Shutdowner,
			dbMigration *versions.DBMigration,
			config core.AppConfig,
		) {
			appLifeCycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Migration App",
						slog.Group("app",
							"version", config.GetAppVersion(),
							"compatible-version", config.GetCompatibleVersion()))

					if err := dbMigration.Migrate(ctx, map[string]migrationCore.Migration{
						migrationCore.DBInitVersion:   versions.NewInitializationMigration(logger),
						config.GetCompatibleVersion(): versions.NewUpgradeMigration(logger),
					}); err != nil {
						return err
					}

					return shutdowner.Shutdown(fx.ExitCode(0))
				},
			})
		}),
	).Run()
}
