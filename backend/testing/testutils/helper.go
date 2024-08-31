package testutils

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing/fstest"
	"wano-island/common/core"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DetectLeakyGoroutines is a helper function that detects and prevents leaked goroutines during testing.
// It uses the gleak library to track goroutine leaks and asserts that no goroutines are leaked.
// This function should be called at the beginning of each test case that involves goroutines.
// It uses ginkgo's DeferCleanup function to ensure that the cleanup code is executed after each test case.
func DetectLeakyGoroutines() {
	// Capture the initial set of goroutines
	nonLeakyGoroutines := Goroutines()

	// Register a cleanup function to assert that no goroutines are leaked
	DeferCleanup(func() {
		// Assert that no goroutines are leaked
		Eventually(Goroutines).ShouldNot(HaveLeaked(nonLeakyGoroutines))
	})
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
		func(postgresCfg *postgres.Config, gormCfg *gorm.Config) {
			postgresCfg.Conn = db
			// Need to disable PrepareStmt to make mocking queries easier
			gormCfg.PrepareStmt = false
		},
	)
	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(func() {
		mock.ExpectClose()

		sqlDB, getDBErr := gormDB.DB()
		Expect(getDBErr).NotTo(HaveOccurred())

		closeDBErr := sqlDB.Close()
		Expect(closeDBErr).NotTo(HaveOccurred())

		Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
	})

	return gormDB, mock
}

// getWorkspaceDir retrieves the workspace directory path by executing the "go env GOWORK" command.
func getWorkspaceDir() string {
	cmd := exec.Command("go", "env", "GOWORK")
	output, _ := cmd.Output()

	return filepath.Dir(string(output))
}

// GetResourceFS creates a virtual file system (fs.FS).
func GetResourceFS() fs.FS {
	localeFile := filepath.Join(getWorkspaceDir(), "console", "resources", "trans", "locale.en.yaml")
	localeEnYamlData, _ := os.ReadFile(localeFile)

	return fstest.MapFS{
		"resources/trans/locale.en.yaml": &fstest.MapFile{
			Data: localeEnYamlData,
		},
	}
}
