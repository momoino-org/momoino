package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	slogGorm "github.com/orandin/slog-gorm"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

// HasCreatedByColumn provides a common column for tracking who created a particular record.
// This struct can be embedded in other models to ensure the "CreatedBy" field is present.
type HasCreatedByColumn struct {
	// The unique identifier (e.g., username or ID) of the user or system that created the record.
	CreatedBy string `gorm:"type:string;size:256;not null"`
}

// HasUpdatedAtColumn is a struct that includes an UpdatedAt field of type time.Time.
// This field is automatically set by GORM when a record is updated.
type HasUpdatedAtColumn struct {
	UpdatedAt time.Time `gorm:"type:time"`
}

type newGormDatabaseParams struct {
	fx.In
	Logger       *slog.Logger
	Config       AppConfig
	Encryptor    Encryptor `name:"aes-gcm"`
	AppLifeCycle fx.Lifecycle
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
	logger *slog.Logger,
	encryptor Encryptor,
	opt func(postgresCfg *postgres.Config, gormCfg *gorm.Config),
) (*gorm.DB, error) {
	postgresCfg := postgres.Config{}
	gormCfg := gorm.Config{
		PrepareStmt: true,
		Logger: slogGorm.New(
			slogGorm.WithHandler(logger.Handler()),
			slogGorm.WithTraceAll(),
		),
	}

	opt(&postgresCfg, &gormCfg)

	schema.RegisterSerializer("encryption", EncryptionSerializer{
		Encryptor: encryptor,
	})

	return gorm.Open(postgres.New(postgresCfg), &gormCfg)
}

// newGormDatabase initializes a new GORM database connection with retry mechanism.
// It also sets up lifecycle hooks for starting and stopping the database connection.
func newGormDatabase(p newGormDatabaseParams) *gorm.DB {
	databaseCfg := p.Config.GetDatabaseConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		databaseCfg.Host,
		databaseCfg.Username,
		databaseCfg.Password,
		databaseCfg.DatabaseName,
		databaseCfg.Port,
	)

	gormDB, openDBConnectionErr := retry.DoWithData(func() (*gorm.DB, error) {
		return OpenDatabase(p.Logger, p.Encryptor, func(postgresCfg *postgres.Config, gormCfg *gorm.Config) {
			postgresCfg.DSN = dsn
		})
	}, retry.Attempts(databaseCfg.MaxAttempts), retry.Delay(1*time.Second))

	p.AppLifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if openDBConnectionErr != nil {
				p.Logger.Error("Cannot connect to database", "details", errors.Unwrap(openDBConnectionErr))

				return errors.Unwrap(openDBConnectionErr)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			p.Logger.Info("Trying to close the database connection")

			sqlDB, getSQLDbErr := gormDB.DB()

			if getSQLDbErr != nil {
				p.Logger.Error("Unable to get *sql.DB instance.", "details", getSQLDbErr)
				return getSQLDbErr
			}

			if closeSQLDbErr := sqlDB.Close(); closeSQLDbErr != nil {
				p.Logger.Error("Unable to close database connection", "details", closeSQLDbErr)
				return closeSQLDbErr
			}

			p.Logger.Info("The database connection closed successfully")

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

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := GetPageSize(r)
		offset := GetOffset(r)

		return db.Offset(offset).Limit(pageSize)
	}
}
