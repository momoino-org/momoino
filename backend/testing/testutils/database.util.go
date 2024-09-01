package testutils

import (
	"context"
	"time"
	"wano-island/common/core"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	_ "github.com/jackc/pgx/v5/stdlib"
	gormPostgres "gorm.io/driver/postgres"
)

// closeGormDB closes the database connection associated with the given gorm.DB instance.
func closeGormDB(gormDB *gorm.DB) {
	sqlDB, getDBErr := gormDB.DB()
	Expect(getDBErr).NotTo(HaveOccurred())

	closeDBErr := sqlDB.Close()
	Expect(closeDBErr).NotTo(HaveOccurred())
}

// CreateDBInstance initializes and returns a gorm.DB instance connected to a PostgreSQL database.
func CreateDBInstance(cfgFn func(postgresCfg *gormPostgres.Config, gormCfg *gorm.Config)) *gorm.DB {
	gormDB, err := core.OpenDatabase(core.NewNoopLogger(), cfgFn)
	Expect(err).ToNot(HaveOccurred())

	DeferCleanup(func() {
		closeGormDB(gormDB)
	})

	return gormDB
}

// CreateTestDBInstance creates a test database instance using sqlmock and gorm.DB.
// It returns a pointer to the gorm.DB instance and a sqlmock.Sqlmock instance for mocking database operations.
// The function also includes a cleanup mechanism using DeferCleanup to ensure that the database connection is closed
// and that all expectations set on the mock are met.
func CreateTestDBInstance() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	Expect(err).NotTo(HaveOccurred())

	gormDB, err := core.OpenDatabase(
		core.NewNoopLogger(),
		func(postgresCfg *gormPostgres.Config, gormCfg *gorm.Config) {
			postgresCfg.Conn = db
			// Need to disable PrepareStmt to make mocking queries easier
			gormCfg.PrepareStmt = false
		},
	)
	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(func() {
		mock.ExpectClose()
		closeGormDB(gormDB)
		Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
	})

	return gormDB, mock
}

// CreatePostgresContainer creates and starts a PostgreSQL container using the testcontainers-go library.
// The function initializes a PostgreSQL container with specific configurations and waits for it to be ready.
func CreatePostgresContainer(ctx context.Context) *postgres.PostgresContainer {
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase("db-migration"),
		postgres.WithUsername("db-migration"),
		postgres.WithPassword("db-migration"),
		postgres.WithSQLDriver("pgx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				//nolint:mnd // No need to fix
				WithOccurrence(2).
				//nolint:mnd // No need to fix
				WithStartupTimeout(5*time.Second)))
	Expect(err).ToNot(HaveOccurred())

	DeferCleanup(func(ctx SpecContext) {
		err := postgresContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	return postgresContainer
}
