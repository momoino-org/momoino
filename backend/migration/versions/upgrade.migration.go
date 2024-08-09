package versions

import (
	"wano-island/common/core"
	migrationCore "wano-island/migration/core"

	"gorm.io/gorm"
)

// upgradeMigration is a struct that implements the core.Migration interface.
// It is responsible for performing database migrations related to upgrading the system.
type upgradeMigration struct {
	logger core.Logger
}

var _ migrationCore.Migration = (*upgradeMigration)(nil)

// NewUpgradeMigration creates a new instance of upgradeMigration.
func NewUpgradeMigration(logger core.Logger) *upgradeMigration {
	return &upgradeMigration{
		logger: logger,
	}
}

// BeforeMigrate is a method that is called before the migration process begins.
func (m *upgradeMigration) BeforeMigrate(_ *gorm.DB) error {
	return nil
}

// Migrate is a method that performs the actual migration operations.
func (m *upgradeMigration) Migrate(_ *gorm.DB) error {
	return nil
}

// AfterMigrate is a method that is called after the migration process is completed.
func (m *upgradeMigration) AfterMigrate(_ *gorm.DB) error {
	return nil
}
