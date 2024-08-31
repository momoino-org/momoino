package testutils

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing/fstest"
	"wano-island/common/core"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CreateTestDBInstance sets up a test database instance using sqlmock and GORM.
// It returns the GORM database instance and the sqlmock instance for mocking database interactions.
func CreateTestDBInstance() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := core.OpenDatabase(
		core.NewNoopLogger(),
		func(postgresCfg *postgres.Config, gormCfg *gorm.Config) {
			postgresCfg.Conn = db
			// Need to disable PrepareStmt to make mocking queries easier
			gormCfg.PrepareStmt = false
		},
	)

	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

// CloseGormDB closes the database connection associated with the provided GORM instance.
func CloseGormDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
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
