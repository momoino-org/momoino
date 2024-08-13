package versions

import (
	"log/slog"
	"wano-island/common/usermgt"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions/initialization"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// initializationMigration is a struct that implements the Migration interface from the core package.
// It is responsible for initializing the database schema.
type initializationMigration struct {
	logger *slog.Logger
}

var _ migrationCore.Migration = (*initializationMigration)(nil)

// NewInitializationMigration returns a new instance of initMigration.
func NewInitializationMigration(logger *slog.Logger) *initializationMigration {
	return &initializationMigration{
		logger: logger,
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
	)
}

// AfterMigrate is a method that is called after the migration process is completed.
func (m *initializationMigration) AfterMigrate(tx *gorm.DB) error {
	encrytedPassword, err := bcrypt.GenerateFromPassword([]byte("Keep!t5ecret"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if result := tx.Create(&usermgt.UserModel{
		Username: "system",
		Email:    "system@internal.com",
		Password: string(encrytedPassword),
	}); result.Error != nil {
		return result.Error
	}

	return nil
}
