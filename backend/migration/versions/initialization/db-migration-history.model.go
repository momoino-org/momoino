package initialization

import "wano-island/common/core"

// DBMigrationModel represents a record of database migration history.
// It is used to track the applied migrations in the application.
type DBMigrationModel struct {
	core.Model
	core.HasCreatedAtColumn

	// Version represents the version of the applied migration.
	Version string `gorm:"type:string;not null;unique;size:256"`
}

// TableName returns the name of the table for the DBMigration struct in the database.
func (DBMigrationModel) TableName() string {
	return "internal.db_migrations"
}
