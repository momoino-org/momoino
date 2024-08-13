package versions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"wano-island/common/core"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions/initialization"

	"github.com/samber/lo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DBMigration is a struct that handles database migration operations.
// It interacts with a GORM database and a logger to perform migration tasks.
type DBMigration struct {
	config core.AppConfig
	db     *gorm.DB
	logger *slog.Logger
}

// newDBMigration creates a new instance of dbMigration.
func NewDBMigration(
	config core.AppConfig,
	db *gorm.DB,
	logger *slog.Logger,
) *DBMigration {
	return &DBMigration{
		config: config,
		db:     db,
		logger: logger,
	}
}

// getDBVersion retrieves the current database version from the db_migrations table.
func (m *DBMigration) getDBVersion(ctx context.Context) (*string, error) {
	migrator := m.db.WithContext(ctx).Migrator()

	if migrator.HasTable(&initialization.DBMigrationModel{}) {
		var dbMigrationModel initialization.DBMigrationModel

		if result := m.db.WithContext(ctx).Order("created_at DESC").First(&dbMigrationModel); result.Error != nil {
			return nil, fmt.Errorf("cannot get the current database version: %w", result.Error)
		}

		return &dbMigrationModel.Version, nil
	}

	return nil, migrationCore.ErrNoDBMigrationTable
}

// getNextMigration determines the next migration to be applied based on the current database version.
//
//nolint:ireturn // We don't know what exactly next migration type will be returned
func (m *DBMigration) getNextMigration(
	ctx context.Context,
	versions map[string]migrationCore.Migration,
) (migrationCore.Migration, error) {
	currentDBVersion, err := m.getDBVersion(ctx)

	if errors.Is(err, migrationCore.ErrNoDBMigrationTable) {
		return versions[migrationCore.DBInitVersion], nil
	}

	if err != nil {
		return nil, err
	}

	if *currentDBVersion == m.config.GetAppVersion() {
		return nil, migrationCore.ErrDBVersionIsUpToDate
	}

	if lo.HasKey(versions, *currentDBVersion) {
		return versions[*currentDBVersion], nil
	}

	return nil, errors.New("there is no migration available for the current database version")
}

// Migrate applies the next migration to the database.
// If the current database version is up to date, logs an info message and returns nil.
func (m *DBMigration) Migrate(
	ctx context.Context,
	versions map[string]migrationCore.Migration,
) error {
	migration, err := m.getNextMigration(ctx, versions)

	if err != nil {
		if errors.Is(err, migrationCore.ErrDBVersionIsUpToDate) {
			m.logger.InfoContext(ctx, "Your db version is up to date. No need to migrate")
			return nil
		}

		return err
	}

	if err = m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if migrateErr := migration.BeforeMigrate(tx); migrateErr != nil {
			m.logger.Error("Failed at [beforeMigrate] step", slog.Any("details", migrateErr))
			return migrateErr
		}

		if migrateErr := migration.Migrate(tx); migrateErr != nil {
			m.logger.Error("Failed at [AutoMigrate] step", slog.Any("details", migrateErr))
			return migrateErr
		}

		if migrateErr := migration.AfterMigrate(tx); migrateErr != nil {
			m.logger.Error("Failed at [AfterAutoMigrate] step", slog.Any("details", migrateErr))
			return migrateErr
		}

		if result := tx.Create(&initialization.DBMigrationModel{
			Version: m.config.GetAppVersion(),
		}); result.Error != nil {
			m.logger.Error("Cannot insert record to public.db_migrations table", slog.Any("details", result.Error))
			return result.Error
		}

		return nil
	}); err != nil {
		m.logger.ErrorContext(ctx, "Failed to migrate database", slog.Any("details", err))
	}

	return err
}

// newDBMigrationModule provides an Fx module for creating a new dbMigration instance.
func NewDBMigrationModule() fx.Option {
	return fx.Module(
		"Database Migration Module",
		fx.Provide(NewDBMigration),
	)
}
