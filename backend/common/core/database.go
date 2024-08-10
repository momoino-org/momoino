package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	slogGorm "github.com/orandin/slog-gorm"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Model is a base struct for all database models. It includes an ID field of type uuid.UUID.
type Model struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid"`
}

// HasCreatedAtColumn is a struct that includes a CreatedAt field of type time.Time.
// This field is automatically set by GORM when a new record is created.
type HasCreatedAtColumn struct {
	CreatedAt time.Time `gorm:"type:time;not null"`
}

// HasUpdatedAtColumn is a struct that includes an UpdatedAt field of type time.Time.
// This field is automatically set by GORM when a record is updated.
type HasUpdatedAtColumn struct {
	UpdatedAt time.Time `gorm:"type:time"`
}

func (u *Model) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		uuidv7, err := uuid.NewV7()

		if err != nil {
			return err
		}

		u.ID = uuidv7
	}

	return nil
}

// OpenDatabase establishes a connection to a PostgreSQL database using GORM.
func OpenDatabase(
	logger Logger,
	cfg postgres.Config,
) (*gorm.DB, error) {
	return gorm.Open(postgres.New(cfg), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
			slogGorm.WithTraceAll(),
			slogGorm.WithContextValue("request-id", RequestIDKey),
		),
	})
}

// newGormDatabase initializes a new GORM database connection with retry mechanism.
// It also sets up lifecycle hooks for starting and stopping the database connection.
func newGormDatabase(
	logger Logger,
	config AppConfig,
	appLifeCycle fx.Lifecycle,
) *gorm.DB {
	databaseCfg := config.GetDatabaseConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		databaseCfg.Host,
		databaseCfg.Username,
		databaseCfg.Password,
		databaseCfg.DatabaseName,
		databaseCfg.Port,
	)

	gormDB, openDBConnectionErr := retry.DoWithData(func() (*gorm.DB, error) {
		return OpenDatabase(logger, postgres.Config{
			DSN: dsn,
		})
	}, retry.Attempts(databaseCfg.MaxAttempts), retry.Delay(1*time.Second))

	appLifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if openDBConnectionErr != nil {
				logger.Error("Cannot connect to database", "details", errors.Unwrap(openDBConnectionErr))

				return errors.Unwrap(openDBConnectionErr)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Trying to close the database connection")

			sqlDB, getSQLDbErr := gormDB.DB()

			if getSQLDbErr != nil {
				logger.Error("Unable to get *sql.DB instance.", "details", getSQLDbErr)
				return getSQLDbErr
			}

			if closeSQLDbErr := sqlDB.Close(); closeSQLDbErr != nil {
				logger.Error("Unable to close database connection", "details", closeSQLDbErr)
				return closeSQLDbErr
			}

			logger.Info("The database connection closed successfully")

			return nil
		},
	})

	return gormDB
}

// NewDatabaseModule provides an fx.Option for registering the newGormDatabase function as a dependency.
func NewDatabaseModule() fx.Option {
	return fx.Module(
		"Database Module",
		fx.Provide(newGormDatabase),
	)
}
