package versions

import (
	"log/slog"
	"wano-island/common/core"
	"wano-island/common/showmgt"
	"wano-island/common/usermgt"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions/initialization"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// initializationMigration is a struct that implements the Migration interface from the core package.
// It is responsible for initializing the database schema.
type initializationMigration struct {
	logger     *slog.Logger
	userSerice usermgt.UserService
}

type InitializationMigrationParams struct {
	Logger     *slog.Logger
	UserSerice usermgt.UserService
}

var _ migrationCore.Migration = (*initializationMigration)(nil)

// NewInitializationMigration returns a new instance of initMigration.
func NewInitializationMigration(p InitializationMigrationParams) *initializationMigration {
	return &initializationMigration{
		logger:     p.Logger,
		userSerice: p.UserSerice,
	}
}

// BeforeMigrate is a method that is called before the migration process begins.
func (m *initializationMigration) BeforeMigrate(tx *gorm.DB) error {
	statements := []string{
		"CREATE SCHEMA IF NOT EXISTS internal",
		"CREATE SCHEMA IF NOT EXISTS public",
		"CREATE EXTENSION IF NOT EXISTS ltree",
	}

	for _, statement := range statements {
		if result := tx.Exec(statement); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Migrate is a method that performs the actual migration process.
func (m *initializationMigration) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(
		&initialization.DBMigrationModel{},
		&usermgt.UserModel{},
		&usermgt.OAuth2ProviderModel{},
		&usermgt.OAuth2UserModel{},
		&showmgt.ShowModel{},
		&showmgt.ShowTranslationModel{},
		&showmgt.SeasonModel{},
		&showmgt.SeasonTranslationModel{},
		&showmgt.EpisodeModel{},
		&showmgt.EpisodeTranslationModel{},
	)
}

// AfterMigrate is a method that is called after the migration process is completed.
func (m *initializationMigration) AfterMigrate(tx *gorm.DB) error {
	encrytedPassword, err := m.userSerice.HashPassword(tx.Statement.Context, "Keep!t5ecret")
	if err != nil {
		return err
	}

	if result := tx.Create(&usermgt.UserModel{
		Username:  "admin",
		FirstName: "Admin",
		Email:     "admin@internal.com",
		Password:  lo.ToPtr(string(encrytedPassword)),
		HasCreatedByColumn: core.HasCreatedByColumn{
			CreatedBy: core.NewSystemUser().GetUsername(),
		},
	}); result.Error != nil {
		return result.Error
	}

	return nil
}
