package core

import "errors"

// ErrNoDBMigrationTable represents an error that occurs when the db_migrations table does not exist.
var ErrNoDBMigrationTable = errors.New("the db_migrations table does not exist")

// ErrDBVersionIsUpToDate represents an error that occurs when the database version is already up to date.
var ErrDBVersionIsUpToDate = errors.New("your database version is up to date")
