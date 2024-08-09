package core

import (
	"gorm.io/gorm"
)

// Migration is an interface that defines methods for database migration operations.
// These methods are used by the application to perform specific actions before, during, and after migrations.
type Migration interface {
	// BeforeMigrate is called before the migration process begins.
	// The function should return an error if any issues occur during the preparation process.
	BeforeMigrate(tx *gorm.DB) error

	// Migrate is called during the migration process.
	// The function should return an error if any issues occur during the migration process.
	Migrate(tx *gorm.DB) error

	// AfterMigrate is called after the migration process is complete.
	// The function should return an error if any issues occur after the migration process.
	AfterMigrate(tx *gorm.DB) error
}
