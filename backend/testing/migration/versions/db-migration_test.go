package versions_test

import (
	"errors"
	"log/slog"
	"regexp"
	"wano-island/common/core"
	migrationCore "wano-island/migration/core"
	"wano-island/migration/versions"
	mockcore "wano-island/testing/mocks/common/core"
	mockmigrationcore "wano-island/testing/mocks/migration/core"
	"wano-island/testing/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var _ = Describe("[migration.versions.db-migration]", func() {
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

	Context("when creating inital database", func() {
		It("should no error when creating inital database", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedConfig.EXPECT().GetAppVersion().Return("1.0.0")
			mockedMigration := mockmigrationcore.NewMockMigration(GinkgoT())
			mockedMigration.EXPECT().BeforeMigrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().Migrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().AfterMigrate(mock.Anything).Return(nil)

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(`INSERT INTO "internal"."db_migrations"`).
				WithArgs(testutils.AnyUUIDArg{}, testutils.AnyTimeArg{}, mockedConfig.GetAppVersion()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: mockedMigration,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should rollback if the BeforeMigrate function failed", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedMigration := mockmigrationcore.NewMockMigration(GinkgoT())
			mockedMigration.EXPECT().BeforeMigrate(mock.Anything).Return(errors.New("something went wrong"))

			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: mockedMigration,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should rollback if the Migrate function failed", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedMigration := mockmigrationcore.NewMockMigration(GinkgoT())
			mockedMigration.EXPECT().BeforeMigrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().Migrate(mock.Anything).Return(errors.New("something went wrong"))

			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: mockedMigration,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should rollback if the AfterMigrate function failed", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedMigration := mockmigrationcore.NewMockMigration(GinkgoT())
			mockedMigration.EXPECT().BeforeMigrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().Migrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().AfterMigrate(mock.Anything).Return(errors.New("something went wrong"))

			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: mockedMigration,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should rollback if cannot insert DB version to the public.db_migrations table", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedConfig.EXPECT().GetAppVersion().Return("1.0.0")
			mockedMigration := mockmigrationcore.NewMockMigration(GinkgoT())
			mockedMigration.EXPECT().BeforeMigrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().Migrate(mock.Anything).Return(nil)
			mockedMigration.EXPECT().AfterMigrate(mock.Anything).Return(nil)

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(`INSERT INTO "internal"."db_migrations"`).
				WithArgs(testutils.AnyUUIDArg{}, testutils.AnyTimeArg{}, mockedConfig.GetAppVersion()).
				WillReturnError(gorm.ErrInvalidValue)
			sqlMock.ExpectRollback()

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{
				migrationCore.DBInitVersion: mockedMigration,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should not run db migration if the database is up to date", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			mockedConfig.EXPECT().GetAppVersion().Return("1.0.0")

			sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM information_schema.tables`)).
				WithArgs("internal", "db_migrations", "BASE TABLE").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "internal"."db_migrations"`)).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("1.0.0"))

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if cannot get the current db version from the database", func(ctx SpecContext) {
			dbMigrator := versions.NewDBMigration(mockedConfig, db, noopLogger)

			sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM information_schema.tables`)).
				WithArgs("internal", "db_migrations", "BASE TABLE").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "internal"."db_migrations"`)).
				WithArgs(1).
				WillReturnError(errors.New("something went wrong"))

			err := dbMigrator.Migrate(ctx, map[string]migrationCore.Migration{})
			Expect(err).To(MatchError(ContainSubstring("cannot get the current database version")))
		})
	})
})
