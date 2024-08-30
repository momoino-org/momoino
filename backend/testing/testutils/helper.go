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
	db, mock, _ := sqlmock.New()
	gormDB, _ := core.OpenDatabase(
		core.NewNoopLogger(),
		postgres.Config{Conn: db},
	)

	return gormDB, mock
}

// CloseGormDB closes the database connection associated with the provided GORM instance.
func CloseGormDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
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
