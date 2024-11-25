package versions_test

import (
	"log/slog"
	"wano-island/common/core"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions"
	"wano-island/migration/versions/initialization"
	mockcore "wano-island/testing/mocks/common/core"
	"wano-island/testing/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("[migration.versions.initialization]", func() {
	Context("when deploy new environment", func() {
		var mockedConfig *mockcore.MockAppConfig
		var db *gorm.DB
		var sqlMock sqlmock.Sqlmock
		var noopLogger *slog.Logger

		BeforeEach(func() {
			testutils.DetectLeakyGoroutines()
			noopLogger = core.NewNoopLogger()
			mockedConfig = mockcore.NewMockAppConfig(GinkgoT())
			db, sqlMock = testutils.CreateTestDBInstance()
		})

		It("should rollback if the pre-process failed", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("CREATE SCHEMA IF NOT EXISTS internal").WillReturnError(gorm.ErrInvalidData)
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: versions.NewInitializationMigration(versions.InitializationMigrationParams{
					Logger: noopLogger,
				}),
			})
			Expect(err).To(HaveOccurred())
		})

		It("should rollback if the Migrate failed", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec("CREATE SCHEMA IF NOT EXISTS internal").WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectExec("CREATE SCHEMA IF NOT EXISTS public").WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectExec("CREATE EXTENSION IF NOT EXISTS ltree").WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectExec(`CREATE TABLE "internal"."db_migrations"`).WillReturnError(gorm.ErrPrimaryKeyRequired)
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: versions.NewInitializationMigration(versions.InitializationMigrationParams{
					Logger: noopLogger,
				}),
			})
			Expect(err).To(HaveOccurred())
		})

		// It("should rollback if the AfterMigrate failed", func(ctx SpecContext) {
		// 	dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

		// 	sqlMock.ExpectBegin()
		// 	sqlMock.ExpectExec("CREATE EXTENSION IF NOT EXISTS ltree").WillReturnResult(sqlmock.NewResult(1, 1))
		// 	sqlMock.ExpectExec(`CREATE TABLE "public"."db_migrations"`).WillReturnResult(sqlmock.NewResult(1, 1))
		// 	sqlMock.ExpectExec(`CREATE TABLE "public"."files"`).WillReturnResult(sqlmock.NewResult(1, 1))
		// 	sqlMock.ExpectExec(`CREATE TABLE "public"."documents"`).WillReturnResult(sqlmock.NewResult(1, 1))
		// 	sqlMock.ExpectRollback()

		// 	err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
		// 		migrationCore.DBInitVersion: versions.NewInitializationMigration(noopLogger),
		// 	})
		// 	Expect(err).To(HaveOccurred())
		// })
	})

	Context("integration test with real database", Ordered, func() {
		var postgresContainer *postgres.PostgresContainer
		var mockedConfig *mockcore.MockAppConfig
		var gormDB *gorm.DB
		var noopLogger *slog.Logger

		BeforeAll(func(ctx SpecContext) {
			postgresContainer = testutils.CreatePostgresContainer(ctx)
			Expect(postgresContainer.Snapshot(ctx, postgres.WithSnapshotName(ctx.SpecReport().FileName()))).ToNot(HaveOccurred())

			gormDB = testutils.CreateDBInstance(func(postgresCfg *gormPostgres.Config, gormCfg *gorm.Config) {
				postgresCfg.DSN = postgresContainer.MustConnectionString(ctx)
			})
		})

		BeforeEach(func() {
			mockedConfig = mockcore.NewMockAppConfig(GinkgoT())
			mockedConfig.EXPECT().GetAppVersion().Return("1.0.0")
		})

		AfterEach(func(ctx SpecContext) {
			err := postgresContainer.Restore(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should no error occurs when creating initialize database", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, gormDB, noopLogger)
			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: versions.NewInitializationMigration(versions.InitializationMigrationParams{
					Logger: noopLogger,
				}),
			})
			Expect(err).ToNot(HaveOccurred())

			var dbMigrationRecords []initialization.DBMigrationModel
			Expect(gormDB.Order("created_at DESC").Find(&dbMigrationRecords).Error).NotTo(HaveOccurred())
			Expect(dbMigrationRecords).To(HaveLen(1))
			Expect(dbMigrationRecords[0].Version).To(Equal("1.0.0"))
		})
	})
})
